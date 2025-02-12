package postgres

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/domain"
)

type Transactions struct {
	db *gorm.DB
}

func NewTransactionsRepository(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

func (r *Transactions) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	r.db.Create(transaction)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return transaction, nil
}

func (r *Transactions) GetPurchasesForUserByUsername(userGUID string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	r.db.Where("receiver_guid = ? OR sender_guid = ?", userGUID, userGUID).Find(&transactions)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return transactions, nil
}
