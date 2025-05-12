package trading

import (
	"errors"
	"fmt"
	"order_go/internal/account"
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
	AddPositionBalanceRatio = 0.98 // 98%
	
	// ClosePositionRatio 平仓时平掉持仓量的比例
	ClosePositionRatio = 0.5 // 50%
	
	// MinPositionRatioThreshold 持仓量占交易对最大交易额度的最小比例阈值，低于该阈值时全部平仓
	MinPositionRatioThreshold = 0.40 // 40%
	
	// MinAddPositionRatioThreshold 加仓时剩余可用资金占交易对最大交易额度的最小比例阈值，低于该阈值时不进行加仓
	MinAddPositionRatioThreshold = 0.1 // 10%
	
	// DefaultMinAmount 默认最小交易量（当无法从数据库获取时使用）
	DefaultMinAmount = 0.5 // Gate.io对HYPE_USDT的最小交易量要求
)

var (
	// ErrNoPositionToClose 没有持仓可平仓错误
	ErrNoPositionToClose = errors.New("没有持仓可平仓")
	
	// ErrSymbolNotFound 交易对不存在错误
	ErrSymbolNotFound = errors.New("交易对在数据库中不存在")
	
	// ErrExceedMaxPositionRatio 超过交易对最大交易额度错误
	ErrExceedMaxPositionRatio = errors.New("超过交易对最大交易额度限制")
	
	// ErrInsufficientAddPositionRatio 加仓时剩余可用资金比例过低错误
	ErrInsufficientAddPositionRatio = errors.New("剩余可用资金比例过低，不进行加仓")
	
	// 注意：ErrInsufficientBalance 错误已在 engine.go 中定义
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
		// 获取交易对的配置信息
		contractCode, err := getFullContractConfig(signal.Symbol)
		if err != nil {
			config.Logger.Errorw("获取交易对配置失败",
				"error", err.Error(),
				"symbol", signal.Symbol,
			)
			return params, err
		}
		
		// 获取账户总价值
		totalValue, err := account.GetTotalValue(ex)
		if err != nil {
			config.Logger.Errorw("获取账户总价值失败",
				"error", err.Error(),
			)
			// 如果无法获取账户总价值，使用默认的平仓比例
			closeAmount := roundAmount(position.Size * ClosePositionRatio, signal.Symbol)
			params.PositionSide = "close"
			params.Amount = closeAmount
			return params, nil
		}
		
		// 计算交易对的最大可用资金
		maxPositionValue := totalValue * contractCode.MaxPositionRatio / 100.0
		
		// 计算当前持仓价值
		currentPositionValue := position.Size * signal.Price
		
		// 计算当前持仓价值占交易对最大可用资金的比例
		currentPositionRatio := currentPositionValue / maxPositionValue
		
		// 根据持仓比例决定平仓策略
		var closeAmount float64
		if currentPositionRatio <= MinPositionRatioThreshold {
			// 如果持仓比例小于或等于最小阈值，全部平仓
			closeAmount = roundAmount(position.Size, signal.Symbol)
			config.Logger.Infow("持仓比例小于最小阈值，全部平仓",
				"symbol", signal.Symbol,
				"position_size", position.Size,
				"current_position_value", currentPositionValue,
				"max_position_value", maxPositionValue,
				"current_position_ratio", currentPositionRatio,
				"min_position_ratio_threshold", MinPositionRatioThreshold,
			)
		} else {
			// 否则使用标准的平仓比例
			closeAmount = roundAmount(position.Size * ClosePositionRatio, signal.Symbol)
			config.Logger.Infow("使用标准平仓比例",
				"symbol", signal.Symbol,
				"position_size", position.Size,
				"close_position_ratio", ClosePositionRatio,
				"close_amount", closeAmount,
				"current_position_ratio", currentPositionRatio,
			)
		}
		
		// 获取交易对的最小交易量
		minAmount := contractCode.MinAmount
		precision := contractCode.AmountPrecision
		
		// 如果计算出的平仓数量小于最小交易量，有两种选择：
		// 1. 使用最小交易量（如果持仓量足够）
		// 2. 全部平仓（如果持仓量小于最小交易量）
		if closeAmount < minAmount {
			if position.Size >= minAmount {
				// 持仓量足够，使用最小交易量
				closeAmount = roundAmount(minAmount, signal.Symbol)
				config.Logger.Infow("计算出的平仓数量小于最小交易量，使用最小交易量",
					"symbol", signal.Symbol,
					"position_size", fmt.Sprintf("%.*f", precision, position.Size),
					"original_close_amount", fmt.Sprintf("%.*f", precision, position.Size * ClosePositionRatio),
					"adjusted_close_amount", fmt.Sprintf("%.*f", precision, closeAmount),
					"min_amount", fmt.Sprintf("%.*f", precision, minAmount),
				)
			} else {
				// 持仓量不足，全部平仓
				closeAmount = roundAmount(position.Size, signal.Symbol)
				config.Logger.Infow("持仓量小于最小交易量，全部平仓",
					"symbol", signal.Symbol,
					"position_size", fmt.Sprintf("%.*f", precision, position.Size),
					"min_amount", fmt.Sprintf("%.*f", precision, minAmount),
				)
			}
		}
		
		config.Logger.Infow("有现有持仓且信号为卖出，执行平仓操作",
			"symbol", signal.Symbol,
			"position_size", fmt.Sprintf("%.*f", precision, position.Size),
			"close_amount", fmt.Sprintf("%.*f", precision, closeAmount),
			"min_amount", fmt.Sprintf("%.*f", precision, minAmount),
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
// 根据交易对最大交易额度计算可下单的数量
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
	
	// 从数据库获取交易对的配置信息
	contractCode, err := getFullContractConfig(symbol)
	if err != nil {
		config.Logger.Errorw("获取交易对配置失败",
			"error", err.Error(),
			"symbol", symbol,
		)
		return 0, err
	}
	
	// 获取报价货币（USDT）
	quoteCurrency := parts[1]
	
	// 获取账户总价值
	totalValue, err := account.GetTotalValue(ex)
	if err != nil {
		config.Logger.Errorw("获取账户总价值失败",
			"error", err.Error(),
		)
		return 0, err
	}
	
	// 计算交易对的最大可用资金（账户总价值 × 最大交易额度比例）
	maxPositionValue := totalValue * contractCode.MaxPositionRatio / 100.0
	
	// 获取交易对当前持仓
	position, err := ex.GetPosition(symbol)
	if err != nil {
		config.Logger.Errorw("获取持仓信息失败",
			"error", err.Error(),
			"symbol", symbol,
		)
		// 如果无法获取持仓信息，假设当前没有持仓
		position = nil
	}
	
	// 计算当前持仓价值
	currentPositionValue := 0.0
	if position != nil && position.Size > 0 {
		currentPositionValue = position.Size * price
	}
	
	// 计算剩余可用资金（最大可用资金 - 当前持仓价值）
	remainingFunds := maxPositionValue - currentPositionValue
	if remainingFunds <= 0 {
		config.Logger.Warnw("已达到或超过交易对最大交易额度限制",
			"symbol", symbol,
			"max_position_value", maxPositionValue,
			"current_position_value", currentPositionValue,
			"max_position_ratio", contractCode.MaxPositionRatio,
		)
		return 0, ErrExceedMaxPositionRatio
	}
	
	// 使用剩余可用资金的一定比例进行开仓
	// 开仓时使用交易对最大可用资金的 InitialOrderBalanceRatio 比例
	desiredFunds := remainingFunds * InitialOrderBalanceRatio
	
	// 获取报价货币（USDT）的可用余额
	available, _, err := ex.GetBalance(quoteCurrency)
	if err != nil {
		config.Logger.Errorw("获取账户余额失败",
			"error", err.Error(),
			"currency", quoteCurrency,
		)
		return 0, err
	}
	
	// 检查账户可用余额是否足够
	if available < desiredFunds {
		config.Logger.Warnw("账户可用余额不足",
			"symbol", symbol,
			"available", available,
			"desired_funds", desiredFunds,
		)
		return 0, ErrInsufficientBalance
	}
	
	config.Logger.Infow("计算开仓数量",
		"symbol", symbol,
		"total_value", totalValue,
		"max_position_ratio", contractCode.MaxPositionRatio,
		"max_position_value", maxPositionValue,
		"current_position_value", currentPositionValue,
		"remaining_funds", remainingFunds,
		"desired_funds", desiredFunds,
		"available", available,
	)
	
	// 计算可买入的数量并根据精度进行四舍五入
	amount := roundAmount(desiredFunds/price, symbol)
	if amount == 0 {
		err := fmt.Errorf("计算的交易数量小于最小交易量 %.5f", contractCode.MinAmount)
		config.Logger.Warnw(err.Error(),
			"symbol", symbol,
			"min_amount", contractCode.MinAmount,
		)
		return 0, err
	}
	
	return amount, nil
}

// calculateAddableAmount 计算可加仓数量
// 根据交易对最大交易额度计算可加仓的数量
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
	
	// 从数据库获取交易对的配置信息
	contractCode, err := getFullContractConfig(symbol)
	if err != nil {
		config.Logger.Errorw("获取交易对配置失败",
			"error", err.Error(),
			"symbol", symbol,
		)
		return 0, err
	}
	
	// 获取报价货币（USDT）
	quoteCurrency := parts[1]
	
	// 获取账户总价值
	totalValue, err := account.GetTotalValue(ex)
	if err != nil {
		config.Logger.Errorw("获取账户总价值失败",
			"error", err.Error(),
		)
		return 0, err
	}
	
	// 计算交易对的最大可用资金（账户总价值 × 最大交易额度比例）
	maxPositionValue := totalValue * contractCode.MaxPositionRatio / 100.0
	
	// 获取交易对当前持仓
	position, err := ex.GetPosition(symbol)
	if err != nil {
		config.Logger.Errorw("获取持仓信息失败",
			"error", err.Error(),
			"symbol", symbol,
		)
		// 如果无法获取持仓信息，假设当前没有持仓
		position = nil
	}
	
	// 计算当前持仓价值
	currentPositionValue := 0.0
	if position != nil && position.Size > 0 {
		currentPositionValue = position.Size * price
	}
	
	// 计算剩余可用资金（最大可用资金 - 当前持仓价值）
	remainingFunds := maxPositionValue - currentPositionValue
	if remainingFunds <= 0 {
		config.Logger.Warnw("已达到或超过交易对最大交易额度限制",
			"symbol", symbol,
			"max_position_value", maxPositionValue,
			"current_position_value", currentPositionValue,
			"max_position_ratio", contractCode.MaxPositionRatio,
		)
		return 0, ErrExceedMaxPositionRatio
	}
	
	// 计算剩余可用资金占交易对最大交易额度的比例
	remainingRatio := remainingFunds / maxPositionValue
	if remainingRatio < MinAddPositionRatioThreshold {
		config.Logger.Warnw("剩余可用资金比例过低，不进行加仓",
			"symbol", symbol,
			"max_position_value", maxPositionValue,
			"current_position_value", currentPositionValue,
			"remaining_funds", remainingFunds,
			"remaining_ratio", remainingRatio,
			"min_add_position_ratio_threshold", MinAddPositionRatioThreshold,
		)
		return 0, ErrInsufficientAddPositionRatio
	}
	
	// 使用剩余可用资金的一定比例进行加仓
	// 加仓时使用交易对最大可用资金的 AddPositionBalanceRatio 比例
	desiredFunds := remainingFunds * AddPositionBalanceRatio
	
	// 获取报价货币（USDT）的可用余额
	available, _, err := ex.GetBalance(quoteCurrency)
	if err != nil {
		config.Logger.Errorw("获取账户余额失败",
			"error", err.Error(),
			"currency", quoteCurrency,
		)
		return 0, err
	}
	
	// 检查账户可用余额是否足够
	if available < desiredFunds {
		config.Logger.Warnw("账户可用余额不足",
			"symbol", symbol,
			"available", available,
			"desired_funds", desiredFunds,
		)
		return 0, ErrInsufficientBalance
	}
	
	config.Logger.Infow("计算加仓数量",
		"symbol", symbol,
		"total_value", totalValue,
		"max_position_ratio", contractCode.MaxPositionRatio,
		"max_position_value", maxPositionValue,
		"current_position_value", currentPositionValue,
		"remaining_funds", remainingFunds,
		"desired_funds", desiredFunds,
		"available", available,
	)
	
	// 计算可买入的数量并根据精度进行四舍五入
	rawAmount := desiredFunds / price
	amount := roundAmount(rawAmount, symbol)
	
	// 记录计算得到的加仓数量
	config.Logger.Infow("计算得到的加仓数量",
		"symbol", symbol,
		"raw_amount", rawAmount,
		"rounded_amount", amount,
		"max_position_ratio", contractCode.MaxPositionRatio,
	)
	
	if amount == 0 {
		err := fmt.Errorf("计算的加仓数量小于最小交易量 %.5f", contractCode.MinAmount)
		config.Logger.Warnw(err.Error(),
			"symbol", symbol,
			"min_amount", contractCode.MinAmount,
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

// getFullContractConfig 获取完整的交易对配置
// 从数据库中读取交易对的完整配置信息
func getFullContractConfig(symbol string) (models.ContractCode, error) {
	// 从数据库中查询交易对配置
	var contractCode models.ContractCode
	result := repository.DB.Where("symbol = ?", symbol).First(&contractCode)
	if result.Error != nil {
		// 如果查询失败，返回错误
		config.Logger.Errorw("交易对在数据库中不存在",
			"error", result.Error.Error(),
			"symbol", symbol,
		)
		return contractCode, ErrSymbolNotFound
	}
	
	return contractCode, nil
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
	
	// 记录原始数量和精度信息
	config.Logger.Debugw("开始处理交易数量精度",
		"symbol", symbol,
		"original_amount", amount,
		"precision", precision,
		"min_amount", minAmount,
	)
	
	// 根据精度进行四舍五入
	factor := float64(1)
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	amount = float64(int(amount*factor)) / factor
	
	// 记录精度处理后的数量
	config.Logger.Debugw("完成数量精度处理",
		"symbol", symbol,
		"rounded_amount", fmt.Sprintf("%.*f", precision, amount),
		"precision_factor", factor,
	)
	
	// 检查是否低于最小交易量
	if amount < minAmount {
		config.Logger.Debugw("计算的数量小于最小交易量",
			"symbol", symbol,
			"calculated_amount", fmt.Sprintf("%.*f", precision, amount),
			"min_amount", minAmount,
		)
		return 0.0
	}
	
	return amount
}