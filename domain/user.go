package domain

import (
	"time"
)

type User struct {
	Username    string    `gorm:"column:username;primaryKey"`
	Password    string    `json:"-" gorm:"column:password;not null"`
	Balance     float64   `json:"balance" gorm:"column:balance;type:decimal(20,8);default:1000"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	AccessToken string    `json:"-" gorm:"-"`
}
