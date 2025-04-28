package queue

import (
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/strategy"
	"order_go/internal/trading"
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
    // 创建一个延迟函数，确保无论如何都会将信号发送到存储队列
    defer func() {
        // 将信号发送到存储队列，无论信号是否有效
        // 使用非阻塞方式发送，避免存储队列满时阻塞处理队列
        select {
        case StoreQueue <- signal:
            // 存储成功时不输出日志，减少日志冗余
        default:
            config.Logger.Warnw("存储队列已满，信号未被存储",
                "symbol", signal.Symbol,
                "action", signal.Action,
            )
        }
    }()
    
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
    
    // 2. 调用交易引擎执行交易逻辑
    engine := trading.GetEngine()
    if err := engine.ProcessSignal(signal); err != nil {
        config.Logger.Errorw("交易执行失败",
            "error", err.Error(),
            "symbol", signal.Symbol,
            "action", signal.Action,
        )
        return
    }
    
    config.Logger.Infow("交易信号处理成功",
        "symbol", signal.Symbol,
        "action", signal.Action,
    )
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
