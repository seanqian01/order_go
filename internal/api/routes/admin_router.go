package routes

import (
	"order_go/internal/api/admin"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册后台管理相关路由
func RegisterAdminRoutes(router *gin.Engine) {
	// 创建API路由组
	apiGroup := router.Group("/api")
	{
		// 信号相关路由
		apiGroup.GET("/signals", admin.GetSignals)
		apiGroup.GET("/signals/:id", admin.GetSignalByID)
		
		// 订单相关路由
		apiGroup.GET("/orders", admin.GetOrders)
		apiGroup.GET("/orders/:id", admin.GetOrderByID)
		
		// 统计数据路由
		apiGroup.GET("/stats", admin.GetStats)
	}
	
	// 保留原有的admin路由组，以便将来可能的扩展
	adminGroup := router.Group("/api/admin")
	{
		// 信号相关路由（与/api路由组相同的处理函数）
		adminGroup.GET("/signals", admin.GetSignals)
		adminGroup.GET("/signals/:id", admin.GetSignalByID)
		
		// 订单相关路由（与/api路由组相同的处理函数）
		adminGroup.GET("/orders", admin.GetOrders)
		adminGroup.GET("/orders/:id", admin.GetOrderByID)
		
		// 统计数据路由（与/api路由组相同的处理函数）
		adminGroup.GET("/stats", admin.GetStats)
	}
}