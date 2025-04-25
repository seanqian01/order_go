package repository

import (
	"context"
	"order_go/internal/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func SaveSignal(ctx context.Context, signal models.TradingSignal) error {
    return DB.WithContext(ctx).Create(&signal).Error
}

// GetSignals 获取信号列表，按照创建时间倒序排列
func GetSignals(ctx context.Context, limit int) ([]models.TradingSignal, error) {
    var signals []models.TradingSignal
    
    if limit <= 0 {
        limit = 100 // 默认限制100条
    }
    
    err := DB.WithContext(ctx).
        Order("created_at DESC").
        Limit(limit).
        Find(&signals).Error
        
    return signals, err
}

// GetSignalsBySymbol 根据交易对获取信号列表，按照创建时间倒序排列
func GetSignalsBySymbol(ctx context.Context, symbol string, limit int) ([]models.TradingSignal, error) {
    var signals []models.TradingSignal
    
    if limit <= 0 {
        limit = 100 // 默认限制100条
    }
    
    err := DB.WithContext(ctx).
        Where("symbol = ?", symbol).
        Order("created_at DESC").
        Limit(limit).
        Find(&signals).Error
        
    return signals, err
}
