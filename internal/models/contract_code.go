package models

import "time"

type ContractCode struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Symbol      string    `json:"symbol" binding:"required"`
    Code        string    `json:"code" binding:"required"`
    ExchangeID  uint      `json:"exchange_id" binding:"required"`
    Status      bool      `json:"status" gorm:"default:true"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ContractCode) TableName() string {
    return "contract_codes"
}
