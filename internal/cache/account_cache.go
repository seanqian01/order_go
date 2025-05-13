package cache

import (
	"fmt"
	"order_go/internal/account"
	"order_go/internal/exchange"
	"order_go/internal/utils/config"
	"sync"
	"time"
)

var (
	// AccountValueCache 存储账户总价值的缓存
	accountValueCache    string
	accountValueCacheMux sync.RWMutex
	lastUpdateTime       time.Time
)

// GetCachedAccountValue 获取缓存的账户总价值
func GetCachedAccountValue() string {
	accountValueCacheMux.RLock()
	defer accountValueCacheMux.RUnlock()
	
	// 如果缓存为空或过期（超过5分钟），返回默认值
	if accountValueCache == "" || time.Since(lastUpdateTime) > 5*time.Minute {
		return "0.00"
	}
	
	return accountValueCache
}

// UpdateAccountValueCache 更新账户总价值缓存
func UpdateAccountValueCache() {
	ex := exchange.NewGateIO()
	accountValue, err := account.GetTotalValue(ex)
	if err != nil {
		config.Logger.Errorw("更新账户总价值缓存失败",
			"error", err.Error(),
		)
		return
	}
	
	formattedValue := fmt.Sprintf("%.2f", accountValue)
	
	accountValueCacheMux.Lock()
	accountValueCache = formattedValue
	lastUpdateTime = time.Now()
	accountValueCacheMux.Unlock()
	
	config.Logger.Infow("账户总价值缓存已更新",
		"value", formattedValue,
		"time", lastUpdateTime.Format(time.RFC3339),
	)
}

// StartAccountValueCacheUpdater 启动账户总价值缓存更新器
func StartAccountValueCacheUpdater() {
	// 立即更新一次缓存
	UpdateAccountValueCache()
	
	// 启动定期更新任务
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			UpdateAccountValueCache()
		}
	}()
	
	config.Logger.Info("账户总价值缓存更新器已启动")
}
