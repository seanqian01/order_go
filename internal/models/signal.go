package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TradingSignal struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SecretKey    string    `json:"secretkey" gorm:"-"` // 不存储到数据库
	Symbol       string    `json:"symbol" binding:"required"`
	Scode        string    `json:"scode" binding:"required"`
	ContractType string    `json:"contractType" binding:"required"`
	Price        float64   `json:"price" binding:"required"`
	Action       string    `json:"action" binding:"required"`
	AlertTitle   string    `json:"alert_title" binding:"required"`
	TimeCircle   string    `json:"time_circle" binding:"required"`
	StrategyID   string    `json:"strategy_id" binding:"required"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (TradingSignal) TableName() string {
	return "trading_signals"
}

func (s *TradingSignal) UnmarshalJSON(data []byte) error {
	type Alias TradingSignal
	aux := &struct {
		Price string `json:"price"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	s.Price, err = strconv.ParseFloat(aux.Price, 64)
	if err != nil {
		return fmt.Errorf("invalid price format: %w", err)
	}
	return nil
}
