package strategy

import (
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"

	"gorm.io/gorm"
)

// TrendingStrategy 趋势策略
type TrendingStrategy struct {
	BaseStrategy // 嵌入基础策略
}

// ValidateSignal 验证信号是否有效
func (s *TrendingStrategy) ValidateSignal(signal models.TradingSignal) (bool, string) {
	// 这里实现趋势策略的验证逻辑
	// 示例：简单验证逻辑，可以根据实际需求进行扩展
	
	// 1. 检查交易对是否支持
	if signal.Symbol == "" {
		return false, "交易对不能为空"
	}
	
	// 2. 检查价格是否合理
	if signal.Price <= 0 {
		return false, "价格必须大于0"
	}
	
	// 3. 检查交易方向是否有效
	if signal.Action != "buy" && signal.Action != "sell" {
		return false, "交易方向必须是buy或sell"
	}
	
	// 4. 检查是否与上一个相同交易对的信号方向相同
	valid, reason := s.checkLastSignalDirection(signal)
	if !valid {
		return false, reason
	}
	
	// 在实际应用中，这里可以添加更复杂的趋势分析逻辑
	// 例如：移动平均线、相对强弱指标(RSI)、MACD等技术指标分析
	
	// 示例中简单返回有效
	return true, ""
}

// checkLastSignalDirection 检查最近的信号方向是否与当前信号相同
// 如果相同则返回无效，不同则返回有效
func (s *TrendingStrategy) checkLastSignalDirection(signal models.TradingSignal) (bool, string) {
	// 查询数据库中最近的一条相同交易对的信号记录
	var lastSignal models.TradingSignal
	
	// 严格使用原始symbol进行查询，不做任何修改
	query := repository.DB.Where("symbol = ?", signal.Symbol)
	
	// 只有当信号ID大于0时才添加ID排除条件
	if signal.ID > 0 {
		query = query.Where("id != ?", signal.ID)
	}
	
	err := query.Order("created_at DESC").First(&lastSignal).Error
	
	// 如果没有找到记录，说明是第一次收到该交易对的信号，直接返回有效
	if err != nil {
		// 检查是否是"记录未找到"错误，这是预期的情况
		if err == gorm.ErrRecordNotFound {
			// 首次接收信号，无需特别输出日志
		} else {
			// 其他错误需要记录
			config.Logger.Warnw("查询历史信号时发生错误",
				"error", err.Error(),
			)
		}
		return true, ""
	}
	
	// 如果最近的信号方向与当前信号相同，则返回无效
	if lastSignal.Action == signal.Action {
		return false, "连续信号方向不能相同，上一个信号已经是 " + signal.Action
	}
	
	// 方向不同，返回有效，无需特别输出日志
	return true, ""
}