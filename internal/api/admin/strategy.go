package admin

import (
	"net/http"
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/strategy"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStrategies 获取策略列表
func GetStrategies(c *gin.Context) {
	var strategies []models.Strategy
	var total int64
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	// 查询总数
	repository.DB.Model(&models.Strategy{}).Count(&total)
	
	// 查询数据，按更新时间倒序排列
	repository.DB.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&strategies)
	
	c.JSON(http.StatusOK, gin.H{
		"items": strategies,
		"total": total,
	})
}

// GetStrategyByID 根据ID获取策略
func GetStrategyByID(c *gin.Context) {
	id := c.Param("id")
	var stra models.Strategy
	
	if err := repository.DB.First(&stra, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "策略不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, stra)
}

// CreateStrategy 创建策略
func CreateStrategy(c *gin.Context) {
	var stra models.Strategy
	if err := c.ShouldBindJSON(&stra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}
	
	// 检查策略代码是否已存在
	var count int64
	repository.DB.Model(&models.Strategy{}).Where("code = ?", stra.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "策略代码已存在",
		})
		return
	}
	
	// 创建策略
	if err := repository.DB.Create(&stra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建策略失败: " + err.Error(),
		})
		return
	}
	
	// 重新初始化策略管理器
	straMgr := strategy.GetManager()
	straMgr.InitStrategies()
	
	c.JSON(http.StatusCreated, stra)
}

// UpdateStrategy 更新策略
func UpdateStrategy(c *gin.Context) {
	id := c.Param("id")
	var stra models.Strategy
	
	// 检查策略是否存在
	if err := repository.DB.First(&stra, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "策略不存在",
		})
		return
	}
	
	// 保存原始创建时间
	createdAt := stra.CreatedAt
	
	// 绑定请求数据
	if err := c.ShouldBindJSON(&stra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}
	
	// 检查策略代码是否与其他策略冲突
	var count int64
	repository.DB.Model(&models.Strategy{}).Where("code = ? AND id != ?", stra.Code, id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "策略代码已被其他策略使用",
		})
		return
	}
	
	// 恢复创建时间，更新修改时间
	stra.CreatedAt = createdAt
	stra.UpdatedAt = time.Now()
	
	// 更新策略
	if err := repository.DB.Save(&stra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新策略失败: " + err.Error(),
		})
		return
	}
	
	// 重新初始化策略管理器
	straMgr := strategy.GetManager()
	straMgr.InitStrategies()
	
	c.JSON(http.StatusOK, stra)
}

// DeleteStrategy 删除策略
func DeleteStrategy(c *gin.Context) {
	id := c.Param("id")
	var stra models.Strategy
	
	// 检查策略是否存在
	if err := repository.DB.First(&stra, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "策略不存在",
		})
		return
	}
	
	// 删除策略
	if err := repository.DB.Delete(&stra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除策略失败: " + err.Error(),
		})
		return
	}
	
	// 重新初始化策略管理器
	straMgr := strategy.GetManager()
	straMgr.InitStrategies()
	
	c.JSON(http.StatusOK, gin.H{
		"message": "策略删除成功",
	})
}
