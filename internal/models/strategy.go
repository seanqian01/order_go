package models

import "time"

type Strategy struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `json:"name" binding:"required"`
    Code      string    `json:"code" binding:"required"`
    Status    bool      `json:"status" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Strategy) TableName() string {
    return "strategies"
}
