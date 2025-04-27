package strategy

import (
	"order_go/internal/models"
)

// Strategy 策略接口定义
type Strategy interface {
	// Init 初始化策略，传入数据库中的策略信息
	Init(dbStrategy models.Strategy)
	
	// GetID 返回策略ID
	GetID() uint
	
	// GetName 返回策略名称
	GetName() string
	
	// IsActive 返回策略是否激活
	IsActive() bool
	
	// ValidateSignal 验证信号是否有效
	ValidateSignal(signal models.TradingSignal) (bool, string)
}

// BaseStrategy 基础策略实现，包含共用字段和方法
type BaseStrategy struct {
	dbStrategy models.Strategy
}

// Init 初始化策略
func (s *BaseStrategy) Init(dbStrategy models.Strategy) {
	s.dbStrategy = dbStrategy
}

// GetID 返回策略ID
func (s *BaseStrategy) GetID() uint {
	return s.dbStrategy.ID
}

// GetName 返回策略名称
func (s *BaseStrategy) GetName() string {
	return s.dbStrategy.Name
}

// IsActive 返回策略是否激活
func (s *BaseStrategy) IsActive() bool {
	return s.dbStrategy.Status
}