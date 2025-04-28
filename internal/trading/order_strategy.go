package trading

import (
	"errors"
	"fmt"
	"order_go/internal/exchange"
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
	"strings"
)

// 下单策略常量
const (
	// InitialOrderBalanceRatio 首次开仓使用可用余额的比例
	InitialOrderBalanceRatio = 0.5 // 50%
	
	// AddPositionBalanceRatio 加仓使用可用余额的比例
	AddPositionBalanceRatio = 0.9 // 90%
	
	// ClosePositionRatio 平仓时平掉持仓量的比例
	ClosePositionRatio = 0.5 // 50%
	
	// DefaultMinAmount 默认最小交易量（当无法从数据库获取时使用）
	DefaultMinAmount = 0.5 // Gate.io对HYPE_USDT的最小交易量要求
	
	// DefaultPrecision 默认精度（当无法从数据库获取时使用）
	DefaultPrecision = 5
)

var (
	// ErrNoPositionToClose 没有持仓可平仓错误
	ErrNoPositionToClose = errors.New("没有持仓可平仓")
	
	// ErrSymbolNotFound 交易对不存在错误
	ErrSymbolNotFound = errors.New("交易对在数据库中不存在")
)

// DetermineSpotOrderStrategy 确定现货交易的下单策略
// 根据持仓情况和账户余额确定最终的现货交易下单参数
// 注意：该函数仅适用于现货交易，合约交易应使用单独的策略函数
func DetermineSpotOrderStrategy(signal models.TradingSignal, ex exchange.Exchange) (models.OrderParams, error) {
	// 初始化下单参数
	params := models.OrderParams{
		Symbol:    signal.Symbol,
		Price:     signal.Price,
		Action:    signal.Action,
		OrderType: "limit", // 默认限价单
	}
	
	// 注意：交易对验证已经在信号处理阶段完成，这里不再重复验证
	
	// 1. 检查当前持仓情况
	position, err := ex.GetPosition(signal.Symbol)
	if err != nil {
		config.Logger.Errorw("获取持仓信息失败",
			"error", err.Error(),
			"symbol", signal.Symbol,
		)
		return params, err
	}
	
	// 2. 根据持仓情况和信号方向确定下单策略
	hasPosition := position != nil && position.Size > 0
	
	if !hasPosition {
		// 没有持仓的情况下，需要根据信号方向决定是否下单
		if signal.Action == "sell" {
			// 如果是卖出信号但没有持仓，直接返回错误
			err := errors.New("当前没有持仓仓位，无法卖出现货")
			config.Logger.Warnw(err.Error(),
				"symbol", signal.Symbol,
				"action", signal.Action,
			)
			return params, err
		}
		
		// 如果是买入信号，正常开仓
		config.Logger.Infow("没有现有持仓，按照信号开仓",
			"symbol", signal.Symbol,
			"action", signal.Action,
		)
		
		// 设置开仓参数
		params.PositionSide = "open"
		amount, err := calculateOrderAmount(signal.Price, signal.Symbol, ex)
		if err != nil {
			config.Logger.Errorw("计算开仓数量失败",
				"error", err.Error(),
				"symbol", signal.Symbol,
				"price", signal.Price,
			)
			return params, err
		}
		
		params.Amount = amount
		return params, nil
	}
	
	// 有持仓的情况
	if signal.Action == "sell" {
		// 信号为卖出，执行平仓操作
		// 只平掉持仓量的一半（根据ClosePositionRatio常量）
		closeAmount := position.Size * ClosePositionRatio
		
		// 获取交易对的最小交易量
		minAmount, _, err := getContractConfig(signal.Symbol)
		if err != nil {
			config.Logger.Errorw("获取交易对配置失败",
				"error", err.Error(),
				"symbol", signal.Symbol,
			)
			return params, err
		}
		
		// 如果计算出的平仓数量小于最小交易量，有两种选择：
		// 1. 使用最小交易量（如果持仓量足够）
		// 2. 全部平仓（如果持仓量小于最小交易量）
		if closeAmount < minAmount {
			if position.Size >= minAmount {
				// 持仓量足够，使用最小交易量
				closeAmount = minAmount
				config.Logger.Infow("计算出的平仓数量小于最小交易量，使用最小交易量",
					"symbol", signal.Symbol,
					"position_size", position.Size,
					"original_close_amount", position.Size * ClosePositionRatio,
					"adjusted_close_amount", closeAmount,
					"min_amount", minAmount,
				)
			} else {
				// 持仓量不足，全部平仓
				closeAmount = position.Size
				config.Logger.Infow("持仓量小于最小交易量，全部平仓",
					"symbol", signal.Symbol,
					"position_size", position.Size,
					"min_amount", minAmount,
				)
			}
		}
		
		config.Logger.Infow("有现有持仓且信号为卖出，执行平仓操作",
			"symbol", signal.Symbol,
			"position_size", position.Size,
			"close_amount", closeAmount,
			"min_amount", minAmount,
		)
		
		params.PositionSide = "close"
		params.Amount = closeAmount
		
		return params, nil
	}
	
	// 信号为买入，且有持仓，执行加仓操作
	config.Logger.Infow("有现有持仓且信号为买入，执行加仓操作",
		"symbol", signal.Symbol,
		"current_position", position.Size,
	)
	
	// 3. 计算可加仓数量
	addableAmount, err := calculateAddableAmount(signal.Symbol, signal.Price, ex)
	if err != nil {
		config.Logger.Errorw("计算可加仓数量失败",
			"error", err.Error(),
			"symbol", signal.Symbol,
		)
		return params, err
	}
	
	if addableAmount <= 0 {
		config.Logger.Warnw("余额不足，无法加仓",
			"symbol", signal.Symbol,
			"price", signal.Price,
		)
		return params, ErrInsufficientBalance
	}
	
	params.PositionSide = "open" // 加仓也是开仓操作
	params.Amount = addableAmount
	
	config.Logger.Infow("确定加仓数量",
		"symbol", signal.Symbol,
		"amount", addableAmount,
	)
	
	return params, nil
}

// calculateOrderAmount 计算下单数量
// 根据价格和可用余额计算可下单的数量
func calculateOrderAmount(price float64, symbol string, ex exchange.Exchange) (float64, error) {
	// 交易对格式为"HYPE_USDT"
	parts := strings.Split(symbol, "_")
	if len(parts) < 2 {
		err := errors.New("无效的交易对格式")
		config.Logger.Errorw(err.Error(),
			"symbol", symbol,
		)
		return 0, err
	}
	
	// 从数据库获取交易对的最小交易量和精度
	minAmount, precision, err := getContractConfig(symbol)
	if err != nil {
		config.Logger.Errorw("获取交易对配置失败",
			"error", err.Error(),
			"symbol", symbol,
		)
		return 0, err
	}
	
	// 获取报价货币（USDT）
	quoteCurrency := parts[1]
	
	// 获取报价货币（USDT）的可用余额
	available, _, err := ex.GetBalance(quoteCurrency)
	if err != nil {
		config.Logger.Errorw("获取账户余额失败",
			"error", err.Error(),
			"currency", quoteCurrency,
		)
		return 0, err
	}
	
	config.Logger.Infow("计算开仓数量",
		"currency", quoteCurrency,
		"available", available,
	)
	
	// 使用可用余额（available）的一定比例进行下单
	usableBalance := available * InitialOrderBalanceRatio
	
	// 计算可买入的数量
	amount := usableBalance / price
	
	// 根据精度进行四舍五入
	factor := float64(1)
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	amount = float64(int(amount*factor)) / factor
	
	// 检查是否低于最小交易量
	if amount < minAmount {
		err := fmt.Errorf("计算的交易数量 %.5f 小于最小交易量 %.5f", amount, minAmount)
		config.Logger.Warnw(err.Error(),
			"symbol", symbol,
			"calculated_amount", amount,
			"min_amount", minAmount,
		)
		return 0, err
	}
	
	return amount, nil
}

// calculateAddableAmount 计算可加仓数量
// 根据当前持仓、价格和可用余额计算可加仓的数量
func calculateAddableAmount(symbol string, price float64, ex exchange.Exchange) (float64, error) {
	// 交易对格式为"HYPE_USDT"
	parts := strings.Split(symbol, "_")
	if len(parts) < 2 {
		err := errors.New("无效的交易对格式")
		config.Logger.Errorw(err.Error(),
			"symbol", symbol,
		)
		return 0, err
	}
	
	// 从数据库获取交易对的最小交易量和精度
	minAmount, precision, err := getContractConfig(symbol)
	if err != nil {
		config.Logger.Errorw("获取交易对配置失败",
			"error", err.Error(),
			"symbol", symbol,
		)
		return 0, err
	}
	
	// 获取报价货币（USDT）
	quoteCurrency := parts[1]
	
	// 获取报价货币（USDT）的可用余额
	available, _, err := ex.GetBalance(quoteCurrency)
	if err != nil {
		config.Logger.Errorw("获取账户余额失败",
			"error", err.Error(),
			"currency", quoteCurrency,
		)
		return 0, err
	}
	
	config.Logger.Infow("计算加仓数量",
		"currency", quoteCurrency,
		"available", available,
	)
	
	// 使用可用余额的一定比例进行加仓
	usableBalance := available * AddPositionBalanceRatio
	
	// 计算可买入的数量
	amount := usableBalance / price
	
	// 根据精度进行四舍五入
	factor := float64(1)
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	amount = float64(int(amount*factor)) / factor
	
	// 检查是否低于最小交易量
	if amount < minAmount {
		err := fmt.Errorf("计算的加仓数量 %.5f 小于最小交易量 %.5f", amount, minAmount)
		config.Logger.Warnw(err.Error(),
			"symbol", symbol,
			"calculated_amount", amount,
			"min_amount", minAmount,
		)
		return 0, err
	}
	
	return amount, nil
}

// 使用engine.go中已定义的getBaseCurrency函数

// getContractConfig 获取交易对配置
// 从数据库中读取交易对的最小交易量和精度等配置
func getContractConfig(symbol string) (minAmount float64, precision int, err error) {
	// 从数据库中查询交易对配置
	var contractCode models.ContractCode
	result := repository.DB.Where("symbol = ?", symbol).First(&contractCode)
	if result.Error != nil {
		// 如果查询失败，返回错误，不使用默认值
		config.Logger.Errorw("交易对在数据库中不存在",
			"error", result.Error.Error(),
			"symbol", symbol,
		)
		return 0, 0, ErrSymbolNotFound
	}
	
	return contractCode.MinAmount, contractCode.AmountPrecision, nil
}

// roundAmount 根据交易对的最小交易量和精度进行调整
func roundAmount(amount float64, symbol string) float64 {
	// 从数据库获取交易对的最小交易量和精度
	minAmount, precision, err := getContractConfig(symbol)
	if err != nil {
		// 如果获取失败，直接返回0，表示无法下单
		config.Logger.Errorw("获取交易对配置失败，无法下单",
			"error", err.Error(),
			"symbol", symbol,
		)
		return 0
	}
	
	// 检查是否低于最小交易量
	if amount < minAmount {
		return 0
	}
	
	// 根据精度进行四舍五入
	factor := float64(1)
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	
	return float64(int(amount*factor)) / factor
}