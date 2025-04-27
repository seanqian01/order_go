package queue

import (
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/strategy"
	"order_go/internal/utils/config"
)

var (
	SignalQueue = make(chan models.TradingSignal, 1000) // 缓冲1000个信号
	StoreQueue  = make(chan models.TradingSignal, 1000) // 存储队列，用于异步保存信号
)

func InitSignalQueue() {
    // 初始化队列消费者
    go func() {
        for signal := range SignalQueue {
            // 调用交易引擎处理信号
            processSignal(signal)
        }
    }()
    
    // 初始化存储队列消费者
    go func() {
        for signal := range StoreQueue {
            // 异步保存信号到数据库
            storeSignal(signal)
        }
    }()
}

// processSignal 处理交易信号的核心逻辑
func processSignal(signal models.TradingSignal) {
    // 1. 调用策略管理器验证信号
    strategyManager := strategy.GetManager()
    valid, reason := strategyManager.ValidateSignal(signal)
    
    if !valid {
        config.Logger.Warnw("信号验证失败: "+reason,
            "symbol", signal.Symbol,
            "action", signal.Action,
        )
        return
    }
    
    // 2. 执行交易逻辑
    // TODO: 添加调用Gate.io API创建订单的逻辑
}

// storeSignal 异步保存信号到数据库
func storeSignal(signal models.TradingSignal) {
    if err := repository.DB.Create(&signal).Error; err != nil {
        config.Logger.Errorw("信号保存失败",
            "error", err.Error(),
            "symbol", signal.Symbol,
            "action", signal.Action,
        )
        return
    }
}
