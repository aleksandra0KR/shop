package repository

import (
	"gorm.io/gorm"
	"shop/domain"
	"shop/internal/repository/postgres"
)

type Repository struct {
	Users        Users
	Merch        Merch
	Purchases    Purchases
	Transactions Transactions
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Users:        postgres.NewUsersRepository(db),
		Purchases:    postgres.NewPurchasesRepository(db),
		Transactions: postgres.NewTransactionsRepository(db),
		Merch:        postgres.NewMerchRepository(db),
	}
}

type Users interface {
	GetUserByUsername(string) (*domain.User, error)
	UpdateUser(*domain.User) error
}

type Merch interface{}

type Purchases interface {
	Create(*domain.Purchase) (*domain.Purchase, error)
	GetPurchasesForUserByUserGUID(string) ([]domain.Purchase, error)
}

type Transactions interface {
	Create(*domain.Transaction) (*domain.Transaction, error)
	GetPurchasesForUserByUsername(string) ([]domain.Transaction, error)
}
