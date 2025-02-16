package postgres

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/domain"
)

type Purchases struct {
	db *gorm.DB
}

func NewPurchasesRepository(db *gorm.DB) *Purchases {
	return &Purchases{db: db}
}

func (r *Purchases) Create(tx *gorm.DB, purchase *domain.Purchase) (*domain.Purchase, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}
	db.Create(purchase)
	if db.Error != nil {
		log.Errorf(db.Error.Error())
		return nil, db.Error
	}
	return purchase, nil

}
func (r *Purchases) GetPurchasesForUserByUsername(username string) ([]domain.Purchase, error) {
	var purchases []domain.Purchase
	db := r.db.Where("user_id= ?", username).Find(&purchases)
	if db.Error != nil {
		log.Errorf(db.Error.Error())
		return nil, db.Error
	}
	return purchases, nil
}
