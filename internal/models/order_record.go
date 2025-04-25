package models

import "time"

type OrderRecord struct {
    ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    SignalID      uint      `json:"signal_id" binding:"required"`
    ExchangeID    uint      `json:"exchange_id" binding:"required"`
    ContractCode  string    `json:"contract_code" binding:"required"`
    OrderType     string    `json:"order_type" binding:"required"`
    Price         float64   `json:"price" binding:"required"`
    Amount        float64   `json:"amount" binding:"required"`
    Status        string    `json:"status" binding:"required"`
    CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (OrderRecord) TableName() string {
    return "order_records"
}
