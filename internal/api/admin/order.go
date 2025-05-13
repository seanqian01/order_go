package admin

import (
	"net/http"
	"order_go/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetOrders 获取订单列表
func GetOrders(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	// 获取筛选参数
	symbol := c.Query("symbol")
	action := c.Query("action")
	status := c.Query("status")
	
	// 计算偏移量
	offset := (page - 1) * limit
	
	// 从数据库获取订单列表
	orders, total, err := repository.GetOrdersPaginated(c, offset, limit, symbol, action, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取订单列表失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"items": orders,
		"total": total,
	})
}

// GetOrderByID 获取订单详情
func GetOrderByID(c *gin.Context) {
	// 获取订单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的订单ID",
		})
		return
	}
	
	// 从数据库获取订单详情
	order, err := repository.GetOrderByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取订单详情失败: " + err.Error(),
		})
		return
	}
	
	if order.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "订单不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, order)
}