package routes

import (
	"order_go/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterSignalRoutes(router *gin.Engine) {
    signalGroup := router.Group("/api")
    {
        signalGroup.POST("/webhook", handlers.HandleSignal)
    }
}
