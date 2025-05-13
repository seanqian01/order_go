package admin

import (
	"net/http"
	"order_go/internal/cache"
	"order_go/internal/repository"

	"github.com/gin-gonic/gin"
)

// GetStats 获取统计数据
func GetStats(c *gin.Context) {
    // 获取信号总数
    signalCount, err := repository.GetSignalCount(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "获取信号总数失败: " + err.Error(),
        })
        return
    }
    
    // 获取订单总数
    orderCount, err := repository.GetOrderCount(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "获取订单总数失败: " + err.Error(),
        })
        return
    }
    
    // 从缓存获取账户总价值
    formattedAccountValue := cache.GetCachedAccountValue()
    
    c.JSON(http.StatusOK, gin.H{
        "signalCount": signalCount,
        "orderCount": orderCount,
        "accountValue": formattedAccountValue,
    })
}

// RefreshAccountValue 手动刷新账户总值
func RefreshAccountValue(c *gin.Context) {
    // 调用缓存更新函数
    cache.UpdateAccountValueCache()
    
    // 获取更新后的账户总值
    formattedAccountValue := cache.GetCachedAccountValue()
    
    c.JSON(http.StatusOK, gin.H{
        "accountValue": formattedAccountValue,
        "message": "账户总值已刷新",
    })
}