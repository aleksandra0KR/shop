package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/domain"
	"shop/internal/repository/postgres"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Repository struct {
	DB           *gorm.DB
	Users        Users
	Merch        Merch
	Purchases    Purchases
	Transactions Transactions
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB:           db,
		Users:        postgres.NewUsersRepository(db),
		Purchases:    postgres.NewPurchasesRepository(db),
		Transactions: postgres.NewTransactionsRepository(db),
		Merch:        postgres.NewMerchRepository(db),
	}
}

type Users interface {
	GetUserByUsername(string) (*domain.User, error)
	UpdateUser(*gorm.DB, *domain.User) error
	CreateUser(*domain.User) (*domain.User, error)
}

type Merch interface {
	GetMerchByName(string) (*domain.Merch, error)
}

type Purchases interface {
	Create(*gorm.DB, *domain.Purchase) (*domain.Purchase, error)
	GetPurchasesForUserByUsername(string) ([]domain.Purchase, error)
}

type Transactions interface {
	Create(*gorm.DB, *domain.Transaction) (*domain.Transaction, error)
	GetTransactionsForUserByUsername(string) ([]domain.Transaction, error)
}

func (r *Repository) CreatePurchase(username string, merchName string) (*domain.Purchase, error) {
	tx := r.DB.Begin()
	user, err := r.Users.GetUserByUsername(username)
	if err != nil {
		tx.Rollback()
		log.Errorf(err.Error())
		return nil, err
	}

	merch, err := r.Merch.GetMerchByName(merchName)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return nil, err
	}

	if user.Balance < (merch.Price) {
		tx.Rollback()
		return nil, errors.New("insufficient money")
	}
	user.Balance -= merch.Price

	purchase := &domain.Purchase{
		User:      *user,
		UserID:    user.Username,
		Merch:     *merch,
		MerchName: merch.Name,
	}
	purchase, err = r.Purchases.Create(tx, purchase)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return nil, err
	}

	err = r.Users.UpdateUser(tx, user)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return purchase, nil
}

func (r *Repository) CreateTransaction(receiverName, senderName string, money float64) (*domain.Transaction, error) {
	tx := r.DB.Begin()
	receiver, err := r.Users.GetUserByUsername(receiverName)
	if err != nil {
		tx.Rollback()
		log.Errorf(err.Error())
		return nil, err
	}

	sender, err := r.Users.GetUserByUsername(senderName)
	if err != nil {
		tx.Rollback()
		log.Errorf(err.Error())
		return nil, err
	}

	if sender.Balance < (money) {
		tx.Rollback()
		return nil, errors.New("insufficient money")
	}

	sender.Balance -= money
	receiver.Balance += money

	transaction := &domain.Transaction{
		Receiver:         *receiver,
		Sender:           *sender,
		MoneyAmount:      money,
		ReceiverUsername: receiver.Username,
		SenderUsername:   sender.Username,
	}

	transaction, err = r.Transactions.Create(tx, transaction)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return nil, err
	}
	err = r.Users.UpdateUser(tx, receiver)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return nil, err
	}
	err = r.Users.UpdateUser(tx, sender)
	if err != nil {
		log.Errorf(err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return transaction, nil
}
