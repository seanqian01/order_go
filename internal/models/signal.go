package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// TradingSignal 交易信号模型
type TradingSignal struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`                // 信号ID
	SecretKey    string    `json:"secretkey" gorm:"-" example:"your-secret-key"`                  // API密钥，不存储到数据库
	Symbol       string    `json:"symbol" binding:"required" example:"BTC_USDT"`                  // 交易对
	Scode        string    `json:"scode" binding:"required" example:"BTC"`                        // 交易对简码
	ContractType int       `json:"contractType" binding:"required" example:"4"`                   // 合约类型: ContractTypeAStock=1(大A股票) ContractTypeFutures=2(商品期货) ContractTypeETF=3(ETF金融指数) ContractTypeCrypto=4(虚拟货币)
	Price        float64   `json:"price" binding:"required" example:"50000.0"`                    // 价格
	Action       string    `json:"action" binding:"required" example:"buy"`                       // 交易动作
	AlertTitle   string    `json:"alert_title" binding:"required" example:"BTC买入信号"`            // 提醒标题
	TimeCircle   string    `json:"time_circle" binding:"required" example:"5m"`                   // 时间周期
	StrategyID   string    `json:"strategy_id" binding:"required" example:"1"`                    // 策略ID
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-04-28T09:00:00+08:00"` // 创建时间
}

// TableName 指定表名
func (TradingSignal) TableName() string {
	return "trading_signals"
}

func (s *TradingSignal) UnmarshalJSON(data []byte) error {
	type Alias TradingSignal
	aux := &struct {
		Price        string `json:"price"`
		ContractType string `json:"contractType"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	// 解析价格
	s.Price, err = strconv.ParseFloat(aux.Price, 64)
	if err != nil {
		return fmt.Errorf("invalid price format: %w", err)
	}
	
	// 解析合约类型
	if aux.ContractType != "" {
		contractType, err := strconv.Atoi(aux.ContractType)
		if err != nil {
			return fmt.Errorf("invalid contractType format: %w", err)
		}
		s.ContractType = contractType
	}
	
	return nil
}
