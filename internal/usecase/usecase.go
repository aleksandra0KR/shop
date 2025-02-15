package usecase

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"shop/domain"
	"shop/internal/repository"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mock.go
type UsecaseImplementation struct {
	Repository *repository.Repository
}

type Usecase interface {
	Auth(username string, password string) (*domain.User, error)
	GetPurchasesForUserByUsername(string) ([]domain.Purchase, error)
	CreateTransaction(string, string, float64) (*domain.Transaction, error)
	CreatePurchase(string, string) (*domain.Purchase, error)
	GetTransactionsForUserByUsername(string) ([]domain.Transaction, error)
}

func NewUsecase(repository *repository.Repository) Usecase {
	return &UsecaseImplementation{Repository: repository}
}

func (r *UsecaseImplementation) Auth(username string, password string) (*domain.User, error) {
	user, err := r.Repository.Users.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (r *UsecaseImplementation) GetPurchasesForUserByUsername(username string) ([]domain.Purchase, error) {
	user, err := r.Repository.Users.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return r.Repository.Purchases.GetPurchasesForUserByUserGUID(user.GUID)
}

func (r *UsecaseImplementation) CreateTransaction(receiver, sender string, money float64) (*domain.Transaction, error) {
	return r.Repository.CreateTransaction(receiver, sender, money)
}

func (r *UsecaseImplementation) CreatePurchase(username string, merchName string) (*domain.Purchase, error) {
	return r.Repository.CreatePurchase(username, merchName)
}

func (r *UsecaseImplementation) GetTransactionsForUserByUsername(username string) ([]domain.Transaction, error) {
	user, err := r.Repository.Users.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return r.Repository.Transactions.GetTransactionsForUserByUserGUID(user.GUID)
}
