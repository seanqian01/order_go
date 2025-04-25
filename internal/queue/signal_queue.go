package queue

import (
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
)

var SignalQueue = make(chan models.TradingSignal, 1000) // 缓冲1000个信号

func InitSignalQueue() {
    // 初始化队列消费者
    go func() {
        for signal := range SignalQueue {
            // 调用交易引擎处理信号
            processSignal(signal)
        }
    }()
}

func processSignal(signal models.TradingSignal) {
    // 保存信号到数据库
    if err := repository.DB.Create(&signal).Error; err != nil {
        config.Logger.Errorw("信号处理失败",
            "error", err.Error(),
            "symbol", signal.Symbol,
            "action", signal.Action,
        )
        return
    }

    config.Logger.Infow("信号处理成功",
        "symbol", signal.Symbol,
        "action", signal.Action,
        "price", signal.Price,
    )
}
