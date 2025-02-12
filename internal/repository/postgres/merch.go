package postgres

import "gorm.io/gorm"

type Merch struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) *Merch {
	return &Merch{db: db}
}
