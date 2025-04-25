package database

import (
	"order_go/internal/models"
	"order_go/internal/utils/config"

	"gorm.io/gorm"
)

// MigrateDB 自动迁移数据库模型
func MigrateDB(db *gorm.DB) {
	config.Logger.Info("开始数据库迁移...")
	
	// 自动迁移模型
	err := db.AutoMigrate(
		&models.TradingSignal{},
		&models.Strategy{},
		&models.TimeCycle{},
		&models.ContractCode{},
		&models.Exchange{},
		&models.User{},
		&models.OrderRecord{},
	)
	if err != nil {
		config.Logger.Errorw("数据库迁移失败", "error", err.Error())
		panic("数据库迁移失败: " + err.Error())
	}
	
	config.Logger.Info("数据库迁移完成")
}