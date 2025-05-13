package admin

import (
	"net/http"
	"order_go/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSignals 获取信号列表
func GetSignals(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	// 获取筛选参数
	symbol := c.Query("symbol")
	action := c.Query("action")
	
	// 计算偏移量
	offset := (page - 1) * limit
	
	// 从数据库获取信号列表
	signals, total, err := repository.GetSignalsPaginated(c, offset, limit, symbol, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取信号列表失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"items": signals,
		"total": total,
	})
}

// GetSignalByID 获取信号详情
func GetSignalByID(c *gin.Context) {
	// 获取信号ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的信号ID",
		})
		return
	}
	
	// 从数据库获取信号详情
	signal, err := repository.GetSignalByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取信号详情失败: " + err.Error(),
		})
		return
	}
	
	if signal.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "信号不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, signal)
}