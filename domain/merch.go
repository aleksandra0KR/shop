package domain

import "github.com/shopspring/decimal"

type Merch struct {
	Guid  string          `json:"guid" gorm:"column:guid;primaryKey"`
	Name  string          `json:"name" gorm:"column:name;not null"`
	Price decimal.Decimal `json:"price" gorm:"column:price;type:decimal(20,8);not null"`
}
