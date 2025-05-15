package trading

import (
	"errors"
	"fmt"
	"order_go/internal/constants"
	"order_go/internal/exchange"
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	engine     *Engine
	engineOnce sync.Once
	
	ErrInvalidContractType = errors.New("无效的合约类型")
	ErrInsufficientBalance = errors.New("账户余额不足")
	ErrOrderFailed         = errors.New("下单失败")
)

// Engine 交易引擎，负责处理交易信号并执行下单操作
type Engine struct {
	exchanges map[string]exchange.Exchange
	monitor   *OrderMonitor
	mutex     sync.RWMutex
}

// GetEngine 获取交易引擎单例
func GetEngine() *Engine {
	engineOnce.Do(func() {
		engine = &Engine{
			exchanges: make(map[string]exchange.Exchange),
			monitor:   GetOrderMonitor(),
		}
		engine.registerExchanges()
	})
	return engine
}

// registerExchanges 注册交易所
func (e *Engine) registerExchanges() {
	// 注册Gate.io交易所
	gateio := exchange.NewGateIO()
	e.exchanges[constants.ExchangeTypeSpot] = gateio
	e.exchanges[constants.ExchangeTypeFutures] = gateio
	
	// 注册交易所到监控器
	e.monitor.RegisterExchange(constants.ExchangeTypeSpot, gateio)
	e.monitor.RegisterExchange(constants.ExchangeTypeFutures, gateio)
	
	// 可以注册更多交易所
	// e.exchanges["binance"] = exchange.NewBinance()
}

// ProcessSignal 处理交易信号，执行下单操作
func (e *Engine) ProcessSignal(signal models.TradingSignal) error {
	// 1. 根据合约类型选择交易所
	ex, exchangeType, err := e.getExchangeByContractType(signal.ContractType)
	if err != nil {
		config.Logger.Errorw("获取交易所失败",
			"error", err.Error(),
			"contract_type", signal.ContractType,
		)
		return err
	}
	
	// 2. 确定下单参数
	orderParams, err := e.determineOrderParams(signal, ex)
	if err != nil {
		config.Logger.Errorw("确定下单参数失败",
			"error", err.Error(),
			"symbol", signal.Symbol,
		)
		return err
	}
	
	// 检查计算出的订单数量是否为0，如果为0则表示不满足最小交易量要求
	if orderParams.Amount <= 0 {
		err := fmt.Errorf("交易数量不足，无法下单")
		config.Logger.Warnw(err.Error(),
			"symbol", signal.Symbol,
			"amount", orderParams.Amount,
		)
		return err
	}
	
	// 3. 创建订单记录
	strategyID, _ := strconv.ParseUint(signal.StrategyID, 10, 64)
	orderRecord := models.OrderRecord{
		StrategyID:   uint(strategyID),
		ExchangeID:   1, // 假设Gate.io的ID为1，实际应从配置或数据库获取
		Symbol:       signal.Symbol,
		ContractType: exchangeType,
		ContractCode: fmt.Sprintf("%d", signal.ContractType), // 存储原始合约类型编码
		OrderType:    orderParams.OrderType,
		Price:        orderParams.Price,
		Amount:       orderParams.Amount,
		Action:       orderParams.Action,
		PositionSide: orderParams.PositionSide,
		Status:       "created",
	}
	
	// 4. 执行下单
	orderReq := &exchange.OrderRequest{
		Symbol:       orderParams.Symbol,
		Price:        orderParams.Price,
		Amount:       orderParams.Amount,
		Side:         orderParams.Action,
		Type:         orderParams.OrderType,
		PositionSide: orderParams.PositionSide,
	}
	
	orderResp, err := ex.CreateOrder(orderReq)
	
	if err != nil {
		config.Logger.Errorw("下单失败",
			"error", err.Error(),
			"symbol", signal.Symbol,
			"action", signal.Action,
		)
		
		// 更新订单状态为失败
		orderRecord.Status = "failed"
		
		// 为失败的订单生成一个唯一的OrderID，避免唯一索引冲突
		// 使用时间戳和随机数组合生成一个临时的OrderID
		timeNow := time.Now().UnixNano()
		orderRecord.OrderID = fmt.Sprintf("failed_%d_%d", timeNow, signal.ID)
		
		if err := repository.DB.Create(&orderRecord).Error; err != nil {
			config.Logger.Errorw("保存失败订单记录失败",
				"error", err.Error(),
			)
		}
		
		return ErrOrderFailed
	}
	
	// 5. 更新订单信息
	orderRecord.OrderID = orderResp.OrderID
	orderRecord.Status = orderResp.Status
	
	// 6. 保存订单信息到数据库
	if err := repository.DB.Create(&orderRecord).Error; err != nil {
		config.Logger.Errorw("保存订单记录失败",
			"error", err.Error(),
			"order_id", orderRecord.OrderID,
		)
		// 继续执行，不返回错误
	}
	
	// 7. 启动订单监控
	e.monitor.StartMonitor(&orderRecord, exchangeType)
	
	return nil
}

// getExchangeByContractType 根据合约类型获取对应的交易所
func (e *Engine) getExchangeByContractType(contractType int) (exchange.Exchange, string, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	// 根据合约类型选择交易所
	var exchangeType string
	
	switch contractType {
	case constants.ContractTypeCrypto: // 虚拟货币
		// 默认使用现货交易
		exchangeType = constants.ExchangeTypeSpot
	default:
		return nil, "", ErrInvalidContractType
	}
	
	ex, ok := e.exchanges[exchangeType]
	if !ok {
		return nil, "", ErrInvalidContractType
	}
	
	return ex, exchangeType, nil
}

// determineOrderParams 确定下单参数
func (e *Engine) determineOrderParams(signal models.TradingSignal, ex exchange.Exchange) (models.OrderParams, error) {
	// 根据合约类型选择不同的下单策略
	// 目前只实现了现货交易策略，合约交易策略待实现
	return DetermineSpotOrderStrategy(signal, ex)
}

// getBaseCurrency 从交易对中提取基础货币
func getBaseCurrency(symbol string) string {
	// 简单实现，假设交易对格式为"BTC_USDT"，基础货币为"USDT"
	// 实际实现应该更复杂，处理不同交易所的不同格式
	parts := strings.Split(symbol, "_")
	if len(parts) > 1 {
		return parts[1]
	}
	return "USDT" // 默认返回USDT
}