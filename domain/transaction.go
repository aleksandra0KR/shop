package domain

import (
	"database/sql/driver"
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Guid              string          `json:"guid" gorm:"column:guid;primaryKey"`
	CreatedAt         time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	ReceiverGUID      string          `json:"receiver_guid" gorm:"column:receiver_guid"`
	Receiver          User            `json:"-" gorm:"foreignKey:ReceiverGUID;references:Guid"`
	SenderGUID        string          `json:"sender_guid" gorm:"column:sender_guid;not null"`
	Sender            User            `json:"-" gorm:"foreignKey:SenderGUID;references:Guid"`
	MoneyAmount       decimal.Decimal `json:"money_amount" gorm:"column:money_amount;type:decimal(20,8);not null"`
	TypeOfTransaction TransactionType `json:"type_of_transaction" gorm:"column:type_of_transaction;type:enum('PURCHASE','TRANSFER'); not null"`
}

type TransactionType string

const (
	TransactionTypePurchase TransactionType = "PURCHASE"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

func (t *TransactionType) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("scan error")
	}
	*t = TransactionType(bytes)
	return nil
}

func (t TransactionType) Value() (driver.Value, error) {
	return string(t), nil
}
