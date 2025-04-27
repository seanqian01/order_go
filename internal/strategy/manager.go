package strategy

import (
	"errors"
	"fmt"
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
	"strconv"
	"sync"
)

var (
	// 单例模式的策略管理器
	manager     *Manager
	managerOnce sync.Once
	
	// 错误定义
	ErrStrategyNotFound = errors.New("策略未找到")
	ErrStrategyInactive = errors.New("策略未激活")
)

// 策略创建器函数类型
type StrategyCreator func() Strategy

// Manager 策略管理器
type Manager struct {
	strategies     map[uint]Strategy
	strategyCreators map[uint]StrategyCreator  // 使用ID作为键
	mutex          sync.RWMutex
}

// GetManager 获取策略管理器单例
func GetManager() *Manager {
	managerOnce.Do(func() {
		manager = &Manager{
			strategies: make(map[uint]Strategy),
			strategyCreators: make(map[uint]StrategyCreator),
		}
		// 注册策略创建器
		manager.registerStrategyCreators()
	})
	return manager
}

// registerStrategyCreators 注册策略创建器
func (m *Manager) registerStrategyCreators() {
	// 注册趋势策略创建器，ID为1
	m.strategyCreators[1] = func() Strategy {
		return &TrendingStrategy{}
	}
	
	// 可以在这里注册更多策略创建器
	// 例如：m.strategyCreators[2] = func() Strategy { return &GridStrategy{} }
}

// initStrategies 初始化策略
func (m *Manager) initStrategies() {
	// 从数据库加载策略
	var dbStrategies []models.Strategy
	if err := repository.DB.Find(&dbStrategies).Error; err != nil {
		config.Logger.Errorw("加载策略失败", "error", err.Error())
		return
	}
	
	// 初始化每个策略
	for _, dbStrategy := range dbStrategies {
		// 根据策略ID获取策略创建器
		creator, ok := m.strategyCreators[dbStrategy.ID]
		if !ok {
			config.Logger.Warnw("未找到策略创建器", "id", dbStrategy.ID)
			continue
		}
		
		// 创建策略实例
		strategy := creator()
		// 初始化策略
		strategy.Init(dbStrategy)
		// 注册策略
		m.strategies[dbStrategy.ID] = strategy
	}
	
	if len(m.strategies) > 0 {
		config.Logger.Infow("策略初始化完成", 
			"count", len(m.strategies),
		)
	} else {
		config.Logger.Warnw("没有找到可用的策略")
	}
}

// InitStrategies 初始化所有策略（公开方法，供外部调用）
func (m *Manager) InitStrategies() {
	m.initStrategies()
}

// GetStrategy 获取策略
func (m *Manager) GetStrategy(id uint) (Strategy, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	strategy, ok := m.strategies[id]
	if !ok {
		return nil, ErrStrategyNotFound
	}
	
	if !strategy.IsActive() {
		return nil, ErrStrategyInactive
	}
	
	return strategy, nil
}

// ValidateSignal 验证信号
func (m *Manager) ValidateSignal(signal models.TradingSignal) (bool, string) {
	// 将字符串ID转换为uint
	strategyID, err := strconv.ParseUint(signal.StrategyID, 10, 32)
	if err != nil {
		return false, fmt.Sprintf("无效的策略ID: %s", signal.StrategyID)
	}
	
	// 获取策略
	strategy, err := m.GetStrategy(uint(strategyID))
	if err != nil {
		return false, err.Error()
	}
	
	// 调用策略验证信号
	return strategy.ValidateSignal(signal)
}