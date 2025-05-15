package routes

import (
	"order_go/internal/api/admin"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册后台管理相关路由
func RegisterAdminRoutes(router *gin.Engine) {
	// 创建普通API路由组，与前端请求路径匹配
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
		apiGroup.POST("/refresh-account", admin.RefreshAccountValue)
		
		// 交易对管理路由
		apiGroup.GET("/contract-codes", admin.GetContractCodes)
		apiGroup.GET("/contract-codes/:id", admin.GetContractCodeByID)
		apiGroup.POST("/contract-codes", admin.CreateContractCode)
		apiGroup.PUT("/contract-codes/:id", admin.UpdateContractCode)
		apiGroup.DELETE("/contract-codes/:id", admin.DeleteContractCode)
		
		// 策略管理路由已移至/api/admin路径下
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
		adminGroup.POST("/refresh-account", admin.RefreshAccountValue)
		
		// 交易对管理路由（与/api路由组相同的处理函数）
		adminGroup.GET("/contract-codes", admin.GetContractCodes)
		adminGroup.GET("/contract-codes/:id", admin.GetContractCodeByID)
		adminGroup.POST("/contract-codes", admin.CreateContractCode)
		adminGroup.PUT("/contract-codes/:id", admin.UpdateContractCode)
		adminGroup.DELETE("/contract-codes/:id", admin.DeleteContractCode)
		
		// 策略管理路由（与/api路由组相同的处理函数）
		adminGroup.GET("/strategies", admin.GetStrategies)
		adminGroup.GET("/strategies/:id", admin.GetStrategyByID)
		adminGroup.POST("/strategies", admin.CreateStrategy)
		adminGroup.PUT("/strategies/:id", admin.UpdateStrategy)
		adminGroup.DELETE("/strategies/:id", admin.DeleteStrategy)
	}
}