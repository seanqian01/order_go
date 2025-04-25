package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Username  string    `json:"username" binding:"required"`
    Password  string    `json:"password" binding:"required"`
    Email     string    `json:"email"`
    Status    bool      `json:"status" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
    return "users"
}
