package handlers

import (
	"net/http"
	"order_go/internal/models"
	"order_go/internal/queue"
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

	// 将信号同时发送到处理队列和存储队列
	signalSent := false
	storeSent := false

	// 尝试发送到处理队列
	select {
	case queue.SignalQueue <- signal:
		signalSent = true
	case <-time.After(100 * time.Millisecond):
		config.Logger.Warn("处理队列已满，信号未被处理",
			"symbol", signal.Symbol,
			"strategy_id", signal.StrategyID)
	}

	// 尝试发送到存储队列
	select {
	case queue.StoreQueue <- signal:
		storeSent = true
	case <-time.After(100 * time.Millisecond):
		config.Logger.Warn("存储队列已满，信号未被存储",
			"symbol", signal.Symbol,
			"strategy_id", signal.StrategyID)
	}

	// 根据发送结果返回不同的状态
	if signalSent && storeSent {
		c.JSON(http.StatusOK, gin.H{
			"status":  "received",
			"symbol":  signal.Symbol,
			"action":  signal.Action,
			"message": "signal processed and stored successfully",
		})
	} else if signalSent {
		c.JSON(http.StatusOK, gin.H{
			"status":  "received",
			"symbol":  signal.Symbol,
			"action":  signal.Action,
			"message": "signal processed but not stored",
		})
	} else if storeSent {
		c.JSON(http.StatusOK, gin.H{
			"status":  "received",
			"symbol":  signal.Symbol,
			"action":  signal.Action,
			"message": "signal stored but not processed",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "failed",
			"symbol":  signal.Symbol,
			"action":  signal.Action,
			"message": "signal neither processed nor stored",
		})
	}
}
