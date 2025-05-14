package admin

import (
	"net/http"
	"order_go/internal/models"
	"order_go/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetContractCodes 获取交易对列表
func GetContractCodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var contractCodes []models.ContractCode
	var total int64

	// 构建查询
	query := repository.DB.Model(&models.ContractCode{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取交易对总数失败: " + err.Error(),
		})
		return
	}

	// 获取分页数据，按更新时间倒序排序
	if err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&contractCodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取交易对列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": contractCodes,
		"total": total,
	})
}

// GetContractCodeByID 根据ID获取交易对详情
func GetContractCodeByID(c *gin.Context) {
	id := c.Param("id")
	var contractCode models.ContractCode

	if err := repository.DB.First(&contractCode, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "交易对不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, contractCode)
}

// CreateContractCode 创建新交易对
func CreateContractCode(c *gin.Context) {
	var contractCode models.ContractCode

	if err := c.ShouldBindJSON(&contractCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}
	
	// 验证最大仓位比例不能为负值
	if contractCode.MaxPositionRatio < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "最大仓位比例不能为负值",
		})
		return
	}

	// 检查交易对是否已存在
	// 使用Count而不是First，避免在没有记录时报错
	var count int64
	if err := repository.DB.Model(&models.ContractCode{}).Where("symbol = ? AND code = ?", contractCode.Symbol, contractCode.Code).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查交易对是否存在时出错: " + err.Error(),
		})
		return
	}
	
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "交易对已存在",
		})
		return
	}

	if err := repository.DB.Create(&contractCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建交易对失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, contractCode)
}

// UpdateContractCode 更新交易对信息
func UpdateContractCode(c *gin.Context) {
	id := c.Param("id")
	var contractCode models.ContractCode

	// 检查交易对是否存在
	if err := repository.DB.First(&contractCode, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "交易对不存在: " + err.Error(),
		})
		return
	}

	// 备份原始数据以便检查唯一性
	originalSymbol := contractCode.Symbol
	originalCode := contractCode.Code

	// 绑定请求数据
	if err := c.ShouldBindJSON(&contractCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}
	
	// 验证最大仓位比例不能为负值
	if contractCode.MaxPositionRatio < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "最大仓位比例不能为负值",
		})
		return
	}

	// 如果Symbol或Code发生变化，检查是否与其他记录冲突
	if (contractCode.Symbol != originalSymbol || contractCode.Code != originalCode) && 
		repository.DB.Where("symbol = ? AND code = ? AND id != ?", 
			contractCode.Symbol, contractCode.Code, contractCode.ID).First(&models.ContractCode{}).RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "交易对已存在",
		})
		return
	}

	// 手动管理更新时间
	
	// 手动设置更新时间为当前时间
	contractCode.UpdatedAt = time.Now()
	
	// 使用Select指定要更新的字段，不包括创建时间
	if err := repository.DB.Model(&contractCode).Select("*").Where("id = ?", contractCode.ID).Updates(map[string]interface{}{
		"symbol":            contractCode.Symbol,
		"code":              contractCode.Code,
		"exchange_id":       contractCode.ExchangeID,
		"min_amount":        contractCode.MinAmount,
		"amount_precision":  contractCode.AmountPrecision,
		"price_precision":   contractCode.PricePrecision,
		"max_position_ratio": contractCode.MaxPositionRatio,
		"status":            contractCode.Status,
		"updated_at":        contractCode.UpdatedAt,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新交易对失败: " + err.Error(),
		})
		return
	}
	
	// 重新查询交易对信息，确保返回最新数据
	if err := repository.DB.First(&contractCode, contractCode.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取更新后的交易对失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, contractCode)
}

// DeleteContractCode 删除交易对
func DeleteContractCode(c *gin.Context) {
	id := c.Param("id")
	var contractCode models.ContractCode

	// 检查交易对是否存在
	if err := repository.DB.First(&contractCode, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "交易对不存在: " + err.Error(),
		})
		return
	}

	// 检查是否有订单或信号引用了该交易对
	var orderCount int64
	if err := repository.DB.Model(&models.OrderRecord{}).Where("contract_code = ?", contractCode.Code).Count(&orderCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查订单关联失败: " + err.Error(),
		})
		return
	}

	if orderCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无法删除交易对，存在关联的订单记录",
		})
		return
	}

	// 删除交易对
	if err := repository.DB.Delete(&contractCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除交易对失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "交易对已删除",
	})
}
