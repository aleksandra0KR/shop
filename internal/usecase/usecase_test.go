package usecase

import (
	"testing"

	"shop/domain"
	"shop/internal/repository"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type (
	MockUsers        struct{ mock.Mock }
	MockPurchases    struct{ mock.Mock }
	MockTransactions struct{ mock.Mock }
)

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

func (m *MockUsers) CreateUser(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	return user, args.Error(0)
}

func (m *MockPurchases) Create(tx *gorm.DB, purchase *domain.Purchase) (*domain.Purchase, error) {
	args := m.Called(tx, purchase)
	if p, ok := args.Get(0).(*domain.Purchase); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPurchases) GetPurchasesForUserByUsername(username string) ([]domain.Purchase, error) {
	args := m.Called(username)
	return args.Get(0).([]domain.Purchase), args.Error(1)
}

func (m *MockTransactions) Create(tx *gorm.DB, transaction *domain.Transaction) (*domain.Transaction, error) {
	args := m.Called(tx, transaction)
	if t, ok := args.Get(0).(*domain.Transaction); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactions) GetTransactionsForUserByUsername(username string) ([]domain.Transaction, error) {
	args := m.Called(username)
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

	mockUsers.On("GetUserByUsername", "user").Return(user, nil)
	authUser, err = usecase.Auth("user", "user2")
	assert.Error(t, err)
	assert.Nil(t, authUser)
	assert.Equal(t, "invalid password", err.Error())

	mockUsers.AssertExpectations(t)
}

func TestGetPurchasesForUserByUsername(t *testing.T) {
	mockUsers := new(MockUsers)
	mockPurchases := new(MockPurchases)
	repo := &repository.Repository{Users: mockUsers, Purchases: mockPurchases}
	usecase := NewUsecase(repo)

	purchases := []domain.Purchase{{UserID: "user"}, {UserID: "user"}}

	mockPurchases.On("GetPurchasesForUserByUsername", "user").Return(purchases, nil)

	result, err := usecase.GetPurchasesForUserByUsername("user")
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockPurchases.On("GetPurchasesForUserByUsername", "user2").Return([]domain.Purchase{}, nil)
	result, err = usecase.GetPurchasesForUserByUsername("user2")
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockUsers.AssertExpectations(t)
	mockPurchases.AssertExpectations(t)
}

func TestGetTransactionsForUserByUsername(t *testing.T) {
	mockUsers := new(MockUsers)
	mockTransactions := new(MockTransactions)
	repo := &repository.Repository{Users: mockUsers, Transactions: mockTransactions}
	usecase := NewUsecase(repo)

	transactions := []domain.Transaction{
		{SenderUsername: "user1", ReceiverUsername: "user2", MoneyAmount: 20},
		{SenderUsername: "user1", ReceiverUsername: "user3", MoneyAmount: 30},
	}

	mockTransactions.On("GetTransactionsForUserByUsername", "user1").Return(transactions, nil)

	result, err := usecase.GetTransactionsForUserByUsername("user1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockTransactions.On("GetTransactionsForUserByUsername", "user8").Return([]domain.Transaction{}, nil)
	result, err = usecase.GetTransactionsForUserByUsername("user8")
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockUsers.AssertExpectations(t)
	mockTransactions.AssertExpectations(t)
}
