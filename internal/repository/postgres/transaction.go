package postgres

import (
	"github.com/google/uuid"
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

func (r *Transactions) Create(tx *gorm.DB, transaction *domain.Transaction) (*domain.Transaction, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}
	transaction.GUID = uuid.New().String()
	db.Create(transaction)
	if db.Error != nil {
		log.Errorf(db.Error.Error())
		return nil, db.Error
	}
	return transaction, nil
}

func (r *Transactions) GetTransactionsForUserByUserGUID(userGUID string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	r.db.Where("receiver_guid = ? OR sender_guid = ?", userGUID, userGUID).Find(&transactions)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return transactions, nil
}
