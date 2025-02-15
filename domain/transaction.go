package domain

import (
	"time"
)

type Transaction struct {
	GUID         string    `json:"guid" gorm:"column:guid;primaryKey"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	ReceiverGUID string    `json:"receiver_guid" gorm:"column:receiver_guid;not null"`
	Receiver     User      `json:"-" gorm:"foreignKey:ReceiverGUID;references:GUID"`
	SenderGUID   string    `json:"sender_guid" gorm:"column:sender_guid;not null"`
	Sender       User      `json:"-" gorm:"foreignKey:SenderGUID;references:GUID"`
	MoneyAmount  float64   `json:"money_amount" gorm:"column:money_amount;type:decimal(20,8);not null"`
}
