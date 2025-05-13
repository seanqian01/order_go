package repository

import (
	"context"
	"order_go/internal/models"
)

// GetSignalsPaginated 分页获取信号列表
func GetSignalsPaginated(ctx context.Context, offset, limit int, symbol, action string) ([]models.TradingSignal, int64, error) {
	var signals []models.TradingSignal
	var total int64
	
	// 构建查询
	query := DB.WithContext(ctx).Model(&models.TradingSignal{})
	
	// 添加筛选条件
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据，按创建时间倒序排序
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&signals).Error; err != nil {
		return nil, 0, err
	}
	
	return signals, total, nil
}

// GetSignalByID 根据ID获取信号详情
func GetSignalByID(ctx context.Context, id uint) (models.TradingSignal, error) {
	var signal models.TradingSignal
	
	if err := DB.WithContext(ctx).First(&signal, id).Error; err != nil {
		return signal, err
	}
	
	return signal, nil
}


// GetSignalCount 获取信号总数
func GetSignalCount(ctx context.Context) (int64, error) {
	var count int64
	if err := DB.WithContext(ctx).Model(&models.TradingSignal{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetOrderCount 获取订单总数
func GetOrderCount(ctx context.Context) (int64, error) {
	var count int64
	if err := DB.WithContext(ctx).Model(&models.OrderRecord{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetOrdersPaginated 分页获取订单列表
func GetOrdersPaginated(ctx context.Context, offset, limit int, symbol, action, status string) ([]models.OrderRecord, int64, error) {
	var orders []models.OrderRecord
	var total int64
	
	// 构建查询
	query := DB.WithContext(ctx).Model(&models.OrderRecord{})
	
	// 添加筛选条件
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据，按更新时间倒序排序
	if err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// GetOrderByID 根据ID获取订单详情
func GetOrderByID(ctx context.Context, id uint) (models.OrderRecord, error) {
	var order models.OrderRecord
	
	if err := DB.WithContext(ctx).First(&order, id).Error; err != nil {
		return order, err
	}
	
	return order, nil
}