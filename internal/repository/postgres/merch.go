package postgres

import (
	"shop/domain"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Merch struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) *Merch {
	return &Merch{db: db}
}

func (r *Merch) GetMerchByName(name string) (*domain.Merch, error) {
	var merch domain.Merch
	r.db.Where("name = ?", name).First(&merch)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return &merch, nil
}
