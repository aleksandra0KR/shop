package domain

import (
	"time"
)

type Transaction struct {
	GUID             string    `json:"guid" gorm:"column:guid;primaryKey"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	ReceiverUsername string    `json:"receiver_username" gorm:"column:receiver_username;not null"`
	Receiver         User      `json:"-" gorm:"foreignKey:ReceiverUsername;references:username"`
	SenderUsername   string    `json:"sender_username" gorm:"column:sender_username;not null"`
	Sender           User      `json:"-" gorm:"foreignKey:SenderUsername;references:username"`
	MoneyAmount      float64   `json:"money_amount" gorm:"column:money_amount;type:decimal(20,8);not null"`
}
