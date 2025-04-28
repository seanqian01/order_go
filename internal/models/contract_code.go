package models

import "time"

type ContractCode struct {
    ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Symbol          string    `json:"symbol" binding:"required"`
    Code            string    `json:"code" binding:"required"`
    ExchangeID      *uint     `json:"exchange_id"`                                // 交易所ID，可为空
    MinAmount       float64   `json:"min_amount" gorm:"default:0.001"`           // 最小交易量
    AmountPrecision int       `json:"amount_precision" gorm:"default:3"`         // 数量精度
    PricePrecision  int       `json:"price_precision" gorm:"default:5"`          // 价格精度
    Status          bool      `json:"status" gorm:"default:true"`
    CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ContractCode) TableName() string {
    return "contract_codes"
}
