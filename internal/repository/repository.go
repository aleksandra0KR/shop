package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/domain"
	"shop/internal/repository/postgres"
)

type Repository struct {
	db           *gorm.DB
	Users        Users
	Merch        Merch
	Purchases    Purchases
	Transactions Transactions
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:           db,
		Users:        postgres.NewUsersRepository(db),
		Purchases:    postgres.NewPurchasesRepository(db),
		Transactions: postgres.NewTransactionsRepository(db),
		Merch:        postgres.NewMerchRepository(db),
	}
}

type Users interface {
	GetUserByUsername(string) (*domain.User, error)
	UpdateUser(*gorm.DB, *domain.User) error
}

type Merch interface {
	GetMerchByName(string) (*domain.Merch, error)
}

type Purchases interface {
	Create(*gorm.DB, *domain.Purchase) (*domain.Purchase, error)
	GetPurchasesForUserByUserGUID(string) ([]domain.Purchase, error)
}

type Transactions interface {
	Create(*gorm.DB, *domain.Transaction) (*domain.Transaction, error)
	GetTransactionsForUserByUserGUID(string) ([]domain.Transaction, error)
}

func (r *Repository) CreatePurchase(username string, merchName string) (error, *domain.Purchase) {
	tx := r.db.Begin()
	user, err := r.Users.GetUserByUsername(username)
	if err != nil {
		tx.Rollback()
		log.Errorf(err.Error())
		return err, nil
	}

	merch, err := r.Merch.GetMerchByName(merchName)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return err, nil
	}

	if user.Balance < (merch.Price) {
		tx.Rollback()
		return errors.New("insufficient money"), nil
	}
	user.Balance -= merch.Price

	purchase := &domain.Purchase{
		User:      *user,
		UserGUID:  user.Guid,
		Merch:     *merch,
		MerchGUID: merch.Guid,
	}
	purchase, err = r.Purchases.Create(tx, purchase)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return err, nil
	}

	err = r.Users.UpdateUser(tx, user)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return err, nil
	}

	tx.Commit()
	return nil, purchase
}

func (r *Repository) CreateTransaction(receiverName, senderName string, money float64) (error, *domain.Transaction) {
	tx := r.db.Begin()
	receiver, err := r.Users.GetUserByUsername(receiverName)
	if err != nil {
		tx.Rollback()
		log.Errorf(err.Error())
		return err, nil
	}

	sender, err := r.Users.GetUserByUsername(senderName)
	if err != nil {
		tx.Rollback()
		log.Errorf(err.Error())
		return err, nil
	}

	if sender.Balance < (money) {
		tx.Rollback()
		return errors.New("insufficient money"), nil
	}

	sender.Balance -= money
	receiver.Balance += money

	transaction := &domain.Transaction{
		Receiver:     *receiver,
		Sender:       *sender,
		MoneyAmount:  money,
		ReceiverGUID: receiver.Guid,
		SenderGUID:   sender.Guid,
	}

	transaction, err = r.Transactions.Create(tx, transaction)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return err, nil
	}
	err = r.Users.UpdateUser(tx, receiver)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return err, nil
	}
	err = r.Users.UpdateUser(tx, sender)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return err, nil
	}
	tx.Commit()
	return nil, transaction
}
