package usecase

import (
	"errors"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"shop/domain"
	"shop/internal/repository"
)

type MockUsers struct{ mock.Mock }
type MockPurchases struct{ mock.Mock }
type MockTransactions struct{ mock.Mock }

func (m *MockUsers) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUsers) UpdateUser(tx *gorm.DB, user *domain.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *MockPurchases) Create(tx *gorm.DB, purchase *domain.Purchase) (*domain.Purchase, error) {
	args := m.Called(tx, purchase)
	if p, ok := args.Get(0).(*domain.Purchase); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPurchases) GetPurchasesForUserByUserGUID(userGUID string) ([]domain.Purchase, error) {
	args := m.Called(userGUID)
	return args.Get(0).([]domain.Purchase), args.Error(1)
}

func (m *MockTransactions) Create(tx *gorm.DB, transaction *domain.Transaction) (*domain.Transaction, error) {
	args := m.Called(tx, transaction)
	if t, ok := args.Get(0).(*domain.Transaction); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactions) GetTransactionsForUserByUserGUID(userGUID string) ([]domain.Transaction, error) {
	args := m.Called(userGUID)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func TestAuth(t *testing.T) {
	mockUsers := new(MockUsers)
	repo := &repository.Repository{Users: mockUsers}
	usecase := NewUsecase(repo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	user := &domain.User{Username: "user", Password: string(hashedPassword)}

	mockUsers.On("GetUserByUsername", "user").Return(user, nil)
	authUser, err := usecase.Auth("user", "user")
	assert.NoError(t, err)
	assert.NotNil(t, authUser)
	assert.Equal(t, "user", authUser.Username)

	mockUsers.On("GetUserByUsername", "user2").Return(nil, nil)
	authUser, err = usecase.Auth("user2", "user2")
	assert.Error(t, err)
	assert.Nil(t, authUser)
	assert.Equal(t, "user not found", err.Error())

	mockUsers.On("GetUserByUsername", "user").Return(user, nil)
	authUser, err = usecase.Auth("user", "user2")
	assert.Error(t, err)
	assert.Nil(t, authUser)
	assert.Equal(t, "invalid password", err.Error())

	mockUsers.AssertExpectations(t)
}

func TestGetPurchasesForUserByUserGUID(t *testing.T) {
	mockUsers := new(MockUsers)
	mockPurchases := new(MockPurchases)
	repo := &repository.Repository{Users: mockUsers, Purchases: mockPurchases}
	usecase := NewUsecase(repo)

	user := &domain.User{GUID: "1", Username: "user"}
	purchases := []domain.Purchase{{UserGUID: "1"}, {UserGUID: "1"}}

	mockUsers.On("GetUserByUsername", "user").Return(user, nil)
	mockPurchases.On("GetPurchasesForUserByUserGUID", "1").Return(purchases, nil)

	result, err := usecase.GetPurchasesForUserByUsername("user")
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockUsers.On("GetUserByUsername", "user2").Return(nil, errors.New("user not found"))
	result, err = usecase.GetPurchasesForUserByUsername("user2")
	assert.Error(t, err)
	assert.Nil(t, result)

	mockUsers.AssertExpectations(t)
	mockPurchases.AssertExpectations(t)
}

func TestGetTransactionsForUserByUserGUID(t *testing.T) {
	mockUsers := new(MockUsers)
	mockTransactions := new(MockTransactions)
	repo := &repository.Repository{Users: mockUsers, Transactions: mockTransactions}
	usecase := NewUsecase(repo)

	user := &domain.User{GUID: "1", Username: "user1"}
	transactions := []domain.Transaction{
		{SenderGUID: "1", ReceiverGUID: "2", MoneyAmount: 20},
		{SenderGUID: "1", ReceiverGUID: "3", MoneyAmount: 30},
	}

	mockUsers.On("GetUserByUsername", "user1").Return(user, nil)
	mockTransactions.On("GetTransactionsForUserByUserGUID", "1").Return(transactions, nil)

	result, err := usecase.GetTransactionsForUserByUsername("user1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockUsers.On("GetUserByUsername", "user8").Return(nil, errors.New("user not found"))
	result, err = usecase.GetTransactionsForUserByUsername("user8")
	assert.Error(t, err)
	assert.Nil(t, result)

	mockUsers.AssertExpectations(t)
	mockTransactions.AssertExpectations(t)
}
