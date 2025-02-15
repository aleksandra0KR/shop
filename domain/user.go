package domain

import (
	"time"
)

type User struct {
	GUID        string    `json:"guid" gorm:"column:guid;primaryKey"`
	Username    string    `json:"username" gorm:"column:username;unique;not null"`
	Password    string    `json:"-" gorm:"column:password;not null"`
	Balance     float64   `json:"balance" gorm:"column:balance;type:decimal(20,8);default:1000"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	AccessToken string    `json:"-" gorm:"-"`
}
