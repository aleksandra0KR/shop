package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"shop/domain"
)

type MockUsers struct{ mock.Mock }
type MockMerch struct{ mock.Mock }
type MockPurchases struct{ mock.Mock }
type MockTransactions struct{ mock.Mock }

func (m *MockUsers) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUsers) UpdateUser(tx *gorm.DB, user *domain.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}
func (m *MockUsers) CreateUser(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockMerch) GetMerchByName(name string) (*domain.Merch, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Merch), args.Error(1)
}

func (m *MockPurchases) Create(tx *gorm.DB, purchase *domain.Purchase) (*domain.Purchase, error) {
	args := m.Called(tx, purchase)
	return args.Get(0).(*domain.Purchase), args.Error(1)
}

func (m *MockPurchases) GetPurchasesForUserByUsername(username string) ([]domain.Purchase, error) {
	args := m.Called(username)
	return args.Get(0).([]domain.Purchase), args.Error(1)
}

func (m *MockTransactions) Create(tx *gorm.DB, transaction *domain.Transaction) (*domain.Transaction, error) {
	args := m.Called(tx, transaction)
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactions) GetTransactionsForUserByUsername(username string) ([]domain.Transaction, error) {
	args := m.Called(username)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func TestCreatePurchase(t *testing.T) {
	mockUsers := new(MockUsers)
	mockMerch := new(MockMerch)
	mockPurchases := new(MockPurchases)
	mockTransactions := new(MockTransactions)
	mockDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	repo := &Repository{
		DB:           mockDB,
		Users:        mockUsers,
		Merch:        mockMerch,
		Purchases:    mockPurchases,
		Transactions: mockTransactions,
	}

	user := &domain.User{Username: "user", Balance: 10000}
	merch := &domain.Merch{Name: "cup", Price: 20}
	purchase := &domain.Purchase{UserID: user.Username, MerchName: merch.Name}

	mockUsers.On("GetUserByUsername", "user").Return(user, nil)
	mockMerch.On("GetMerchByName", "cup").Return(merch, nil)
	mockPurchases.On("Create", mock.Anything, mock.Anything).Return(purchase, nil)
	mockUsers.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	result, err := repo.CreatePurchase("user", "cup")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, purchase, result)

	user.Balance = 10
	result, err = repo.CreatePurchase("user", "cup")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "insufficient money", err.Error())
}

func TestCreateTransaction(t *testing.T) {
	mockUsers := new(MockUsers)
	mockMerch := new(MockMerch)
	mockPurchases := new(MockPurchases)
	mockTransactions := new(MockTransactions)
	mockDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	repo := &Repository{
		DB:           mockDB,
		Users:        mockUsers,
		Merch:        mockMerch,
		Purchases:    mockPurchases,
		Transactions: mockTransactions,
	}

	sender := &domain.User{Username: "user1", Balance: 1000}
	receiver := &domain.User{Username: "user2", Balance: 1000}
	transaction := &domain.Transaction{SenderUsername: sender.Username, ReceiverUsername: receiver.Username, MoneyAmount: 20}

	mockUsers.On("GetUserByUsername", "user1").Return(sender, nil)
	mockUsers.On("GetUserByUsername", "user2").Return(receiver, nil)
	mockTransactions.On("Create", mock.Anything, mock.Anything).Return(transaction, nil)
	mockUsers.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	result, err := repo.CreateTransaction("user2", "user1", 30)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transaction, result)

	sender.Balance = 10
	result, err = repo.CreateTransaction("user2", "user1", 20)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "insufficient money", err.Error())
}
