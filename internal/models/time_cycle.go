package models

import "time"

type TimeCycle struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `json:"name" binding:"required"`
    Code      string    `json:"code" binding:"required"`
    Minutes   int       `json:"minutes" binding:"required"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (TimeCycle) TableName() string {
    return "time_cycles"
}
