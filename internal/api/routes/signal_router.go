package routes

import (
	"order_go/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterSignalRoutes 注册信号相关路由
// @title           Order Go API
// @version         1.0
// @description     交易信号处理系统 API
func RegisterSignalRoutes(router *gin.Engine) {
    // 创建API路由组
    signalGroup := router.Group("/api")
    {
        // 信号处理接口
        signalGroup.POST("/webhook", handlers.HandleSignal)
    }
}
