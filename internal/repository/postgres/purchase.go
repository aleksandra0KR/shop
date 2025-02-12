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

func (r *Purchases) Create(purchase *domain.Purchase) (*domain.Purchase, error) {
	r.db.Create(purchase)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return purchase, nil
}
func (r *Purchases) GetPurchasesForUserByUserGUID(userGUID string) ([]domain.Purchase, error) {
	var purchases []domain.Purchase
	err := r.db.Where("user_guid = ?", userGUID).Find(&purchases)
	if err != nil {
		log.Errorf(err.Error.Error())
		return nil, err.Error
	}
	return purchases, nil
}
