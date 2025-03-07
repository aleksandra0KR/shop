// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	domain "shop/domain"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// GetUserByUsername mocks_repo base method.
func (m *MockUsers) GetUserByUsername(arg0 string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUsersMockRecorder) GetUserByUsername(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUsers)(nil).GetUserByUsername), arg0)
}

// UpdateUser mocks_repo base method.
func (m *MockUsers) UpdateUser(arg0 *gorm.DB, arg1 *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUsersMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUsers)(nil).UpdateUser), arg0, arg1)
}

// MockMerch is a mock of Merch interface.
type MockMerch struct {
	ctrl     *gomock.Controller
	recorder *MockMerchMockRecorder
}

// MockMerchMockRecorder is the mock recorder for MockMerch.
type MockMerchMockRecorder struct {
	mock *MockMerch
}

// NewMockMerch creates a new mock instance.
func NewMockMerch(ctrl *gomock.Controller) *MockMerch {
	mock := &MockMerch{ctrl: ctrl}
	mock.recorder = &MockMerchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerch) EXPECT() *MockMerchMockRecorder {
	return m.recorder
}

// GetMerchByName mocks_repo base method.
func (m *MockMerch) GetMerchByName(arg0 string) (*domain.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerchByName", arg0)
	ret0, _ := ret[0].(*domain.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerchByName indicates an expected call of GetMerchByName.
func (mr *MockMerchMockRecorder) GetMerchByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerchByName", reflect.TypeOf((*MockMerch)(nil).GetMerchByName), arg0)
}

// MockPurchases is a mock of Purchases interface.
type MockPurchases struct {
	ctrl     *gomock.Controller
	recorder *MockPurchasesMockRecorder
}

// MockPurchasesMockRecorder is the mock recorder for MockPurchases.
type MockPurchasesMockRecorder struct {
	mock *MockPurchases
}

// NewMockPurchases creates a new mock instance.
func NewMockPurchases(ctrl *gomock.Controller) *MockPurchases {
	mock := &MockPurchases{ctrl: ctrl}
	mock.recorder = &MockPurchasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPurchases) EXPECT() *MockPurchasesMockRecorder {
	return m.recorder
}

// Create mocks_repo base method.
func (m *MockPurchases) Create(arg0 *gorm.DB, arg1 *domain.Purchase) (*domain.Purchase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*domain.Purchase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPurchasesMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPurchases)(nil).Create), arg0, arg1)
}

// GetPurchasesForUserByUserGUID mocks_repo base method.
func (m *MockPurchases) GetPurchasesForUserByUserGUID(arg0 string) ([]domain.Purchase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPurchasesForUserByUserGUID", arg0)
	ret0, _ := ret[0].([]domain.Purchase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPurchasesForUserByUserGUID indicates an expected call of GetPurchasesForUserByUserGUID.
func (mr *MockPurchasesMockRecorder) GetPurchasesForUserByUserGUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPurchasesForUserByUserGUID", reflect.TypeOf((*MockPurchases)(nil).GetPurchasesForUserByUserGUID), arg0)
}

// MockTransactions is a mock of Transactions interface.
type MockTransactions struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionsMockRecorder
}

// MockTransactionsMockRecorder is the mock recorder for MockTransactions.
type MockTransactionsMockRecorder struct {
	mock *MockTransactions
}

// NewMockTransactions creates a new mock instance.
func NewMockTransactions(ctrl *gomock.Controller) *MockTransactions {
	mock := &MockTransactions{ctrl: ctrl}
	mock.recorder = &MockTransactionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactions) EXPECT() *MockTransactionsMockRecorder {
	return m.recorder
}

// Create mocks_repo base method.
func (m *MockTransactions) Create(arg0 *gorm.DB, arg1 *domain.Transaction) (*domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionsMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactions)(nil).Create), arg0, arg1)
}

// GetTransactionsForUserByUserGUID mocks_repo base method.
func (m *MockTransactions) GetTransactionsForUserByUserGUID(arg0 string) ([]domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsForUserByUserGUID", arg0)
	ret0, _ := ret[0].([]domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsForUserByUserGUID indicates an expected call of GetTransactionsForUserByUserGUID.
func (mr *MockTransactionsMockRecorder) GetTransactionsForUserByUserGUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsForUserByUserGUID", reflect.TypeOf((*MockTransactions)(nil).GetTransactionsForUserByUserGUID), arg0)
}
