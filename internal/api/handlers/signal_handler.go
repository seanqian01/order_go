package handlers

import (
	"context"
	"net/http"
	"order_go/internal/models"
	"order_go/internal/queue"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleSignal(c *gin.Context) {
	var signal models.TradingSignal

	if err := c.ShouldBindJSON(&signal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证密钥
	if signal.SecretKey != config.AppConfig.Server.SecretKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid secret key"})
		return
	}

	// 记录接收到的信号
	config.Logger.Infow("接收到交易信号",
		"symbol", signal.Symbol,
		"action", signal.Action,
		"price", signal.Price,
		"time_circle", signal.TimeCircle,
	)

	// 将信号加入处理队列
	select {
	case queue.SignalQueue <- signal:
		// 队列接收成功，启动异步存储
		go func(s models.TradingSignal) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			if err := repository.SaveSignal(ctx, s); err != nil {
				// 记录存储失败日志
				config.Logger.Errorw("failed to save signal",
					"error", err.Error(),
					"symbol", s.Symbol,
					"strategy_id", s.StrategyID)
			}
		}(signal)
	case <-time.After(100 * time.Millisecond):
		config.Logger.Warn("signal queue is full, dropping signal",
			"symbol", signal.Symbol,
			"strategy_id", signal.StrategyID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "received",
		"symbol":  signal.Symbol,
		"action":  signal.Action,
		"message": "signal processed successfully",
	})
}
