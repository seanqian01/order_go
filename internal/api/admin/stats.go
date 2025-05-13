package admin

import (
	"fmt"
	"net/http"
	"order_go/internal/account"
	"order_go/internal/exchange"
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
    
    // 获取账户总价值
    ex := exchange.NewGateIO()
    accountValue, err := account.GetTotalValue(ex)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "获取账户总价值失败: " + err.Error(),
        })
        return
    }
    
    // 格式化账户总价值，保留2位小数
    formattedAccountValue := fmt.Sprintf("%.2f", accountValue)
    
    c.JSON(http.StatusOK, gin.H{
        "signalCount": signalCount,
        "orderCount": orderCount,
        "accountValue": formattedAccountValue,
    })
}