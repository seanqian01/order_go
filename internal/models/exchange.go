package models

import "time"

type Exchange struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `json:"name" binding:"required"`
    Code      string    `json:"code" binding:"required"`
    ApiKey    string    `json:"api_key"`
    ApiSecret string    `json:"api_secret"`
    Status    bool      `json:"status" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Exchange) TableName() string {
    return "exchanges"
}
