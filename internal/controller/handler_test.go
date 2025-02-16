package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shop/domain"
	"shop/internal/controller/middleware"
	mockusecase "shop/internal/usecase/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	usecaseMock := new(mockusecase.MockUsecase)
	handler := NewHandler(usecaseMock)
	return handler.Handle().(*gin.Engine)
}

func TestUnauthorizedAccess(t *testing.T) {
	testTable := []struct {
		name                 string
		method               string
		path                 string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "UnauthorizedAccess_InfoHandler",
			method:               http.MethodGet,
			path:                 "/api/info",
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"Missing accessToken"}`,
		},
		{
			name:                 "UnauthorizedAccess_SendCoinHandle",
			method:               http.MethodPost,
			path:                 "/api/sendCoin",
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"Missing accessToken"}`,
		},
		{
			name:                 "UnauthorizedAccess_BuyItemHandler(",
			method:               http.MethodPost,
			path:                 "/api/buy/1",
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"Missing accessToken"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			router := setupRouter()
			req := httptest.NewRequest(test.method, test.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			assert.Equal(t, test.expectedStatusCode, rec.Code)
			assert.Equal(t, test.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestInfoHandler_Success_EmptyBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "test")

	mockUsecase.EXPECT().GetPurchasesForUserByUserGUID("test").Return([]domain.Purchase{}, nil)
	mockUsecase.EXPECT().GetTransactionsForUserByUserGUID("test").Return([]domain.Transaction{}, nil)
	expectedResponseBody := `{"purchases":[],"transactions":[]}`
	h.InfoHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestInfoHandler_Success_WithBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "test")

	purchase := domain.Purchase{GUID: "1", UserID: "user1", MerchName: "socks", CreatedAt: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)}
	transaction := domain.Transaction{GUID: "1", ReceiverUsername: "user2", SenderUsername: "user1", MoneyAmount: 100, CreatedAt: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)}

	mockUsecase.EXPECT().GetPurchasesForUserByUserGUID("test").Return([]domain.Purchase{purchase}, nil)
	mockUsecase.EXPECT().GetTransactionsForUserByUserGUID("test").Return([]domain.Transaction{transaction}, nil)
	expectedResponseBody := `{"purchases":[{"guid":"1","user_id":"user1","merch_name":"socks","created_at":"0001-01-01T00:00:00Z"}],"transactions":[{"guid":"1","created_at":"0001-01-01T00:00:00Z","receiver_username":"user2","sender_username":"user1","money_amount":100}]}`
	h.InfoHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestInfoHandler_InternalServerError_Purchases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "test")

	mockUsecase.EXPECT().GetPurchasesForUserByUserGUID("test").Return([]domain.Purchase{}, errors.New("db error"))
	expectedResponseBody := `{"error":"db error"}`
	h.InfoHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestInfoHandler_InternalServerError_Transactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "test")

	mockUsecase.EXPECT().GetPurchasesForUserByUserGUID("test").Return([]domain.Purchase{}, nil)
	mockUsecase.EXPECT().GetTransactionsForUserByUserGUID("test").Return([]domain.Transaction{}, errors.New("db error"))
	expectedResponseBody := `{"error":"db error"}`
	h.InfoHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestSendCoinHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(`{"receiver_username":"receiver","amount":10.0}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "sender")

	transaction := domain.Transaction{GUID: "1", ReceiverUsername: "user2", SenderUsername: "user1", MoneyAmount: 10.0, CreatedAt: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)}
	mockUsecase.EXPECT().CreateTransaction("receiver", "sender", 10.0).Return(&transaction, nil)
	expectedResponseBody := `{"guid":"1","created_at":"0001-01-01T00:00:00Z","receiver_username":"user2","sender_username":"user1","money_amount":10}`

	h.SendCoinHandler(c)
	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSendCoinHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(`{"receiver_username":"receiver","amount":10.0}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "sender")

	mockUsecase.EXPECT().CreateTransaction("receiver", "sender", 10.0).Return(nil, errors.New("db error"))
	expectedResponseBody := `{"error":"db error"}`

	h.SendCoinHandler(c)
	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSendCoinHandler_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(`{"receiver_username":"receiver","amount":10.0}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "sender")

	mockUsecase.EXPECT().CreateTransaction("receiver", "sender", 10.0).Return(nil, errors.New("insufficient money"))
	expectedResponseBody := `{"error":"insufficient money"}`

	h.SendCoinHandler(c)
	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendCoinHandler_BadRequest_MissingFields(t *testing.T) {
	router := setupRouter()
	user := domain.User{Username: "test", Password: "test"}
	validToken, err := middleware.JWT{}.GenerateToken(&user)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBufferString(``))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: validToken})

	router.ServeHTTP(w, req)
	expectedResponseBody := `{"error":"Invalid request"}`

	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendCoinHandler_BadRequest_WrongFieldNames(t *testing.T) {
	router := setupRouter()
	user := domain.User{Username: "test", Password: "test"}
	validToken, err := middleware.JWT{}.GenerateToken(&user)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBufferString(`{"receiver": "user2", "value": 10}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: validToken})

	router.ServeHTTP(w, req)
	expectedResponseBody := `{"error":"Missing or invalid fields"}`

	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyItemHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "item", Value: "sock"})
	c.Set("username", "buyer")

	purchase := domain.Purchase{GUID: "1", UserID: "1", MerchName: "1", CreatedAt: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)}
	mockUsecase.EXPECT().CreatePurchase("buyer", "sock").Return(&purchase, nil)
	expectedResponseBody := `{"guid":"1","user_id":"1","merch_name":"1","created_at":"0001-01-01T00:00:00Z"}`

	h.BuyItemHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestBuyItemHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "item", Value: "sock"})
	c.Set("username", "buyer")

	mockUsecase.EXPECT().CreatePurchase("buyer", "sock").Return(nil, errors.New("db error"))
	expectedResponseBody := `{"error":"db error"}`

	h.BuyItemHandler(c)

	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBuyItemHandler_BadRequest_InsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/buy/socks", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "buyer")
	c.Params = gin.Params{
		{Key: "item", Value: "socks"},
	}

	mockUsecase.EXPECT().CreatePurchase("buyer", "socks").Return(nil, errors.New("insufficient money"))
	expectedResponseBody := `{"error":"insufficient money"}`

	h.BuyItemHandler(c)
	assert.Equal(t, expectedResponseBody, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockusecase.NewMockUsecase(ctrl)
	h := NewHandler(mockUsecase)

	t.Run("Invalid JSON Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(`{"username": "user1", "password": "user1"`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AuthHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid request"}`, w.Body.String())
	})

	t.Run("Missing or Invalid Fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(`{"username": "", "password": "user1"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		h.AuthHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Missing or invalid fields"}`, w.Body.String())

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(`{"username": "user1", "password": ""}`))
		c.Request.Header.Set("Content-Type", "application/json")
		h.AuthHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Missing or invalid fields"}`, w.Body.String())
	})

	t.Run("Authentication Failure", func(t *testing.T) {
		mockUsecase.EXPECT().Auth("user1", "user1user1").Return(nil, errors.New("invalid password"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(`{"username": "user1", "password": "user1user1"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		h.AuthHandler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "invalid password"}`, w.Body.String())
	})

	t.Run("Successful Authentication", func(t *testing.T) {
		mockUser := &domain.User{Username: "user1", Password: "user1"}
		mockUsecase.EXPECT().Auth("user1", "user1").Return(mockUser, nil)

		mockToken, err := middleware.JWT{}.GenerateToken(mockUser)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(`{"username": "user1", "password": "user1"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		h.AuthHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"response": {"accessToken": "`+mockToken+`"}}`, w.Body.String())
	})
}
