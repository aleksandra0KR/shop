//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shop/domain"
	"shop/internal/controller"
	"shop/internal/controller/middleware"
	"shop/internal/repository"
	"shop/internal/usecase"
	hash "shop/pkg"
	"shop/pkg/database"
	"shop/pkg/logger"
	"testing"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() (http.Handler, usecase.Usecase, *gorm.DB) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("error loading .env file")
	}

	db := database.InitializeDBPostgres(3, 10)
	clearDatabase(db.GetDB())
	db.Seed()
	logger.InitLogger()

	repository := repository.NewRepository(db.GetDB())
	usecase := usecase.NewUsecase(repository)
	handler := controller.NewHandler(usecase)
	router := handler.Handle()

	return router, usecase, db.GetDB()
}

func performAuthRequest(t *testing.T, router http.Handler, username, password string) string {
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(reqBodyJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, resBody["response"], "accessToken")

	token := resBody["response"].(map[string]interface{})["accessToken"].(string)
	return token
}

func TestAuthHandlerIntegration(t *testing.T) {
	router, _, db := setupTestDB()
	defer clearDatabase(db)

	user := domain.User{Username: "testuser", Password: hash.HashPassword("user1"), Balance: 100}
	db.Create(&user)

	token := performAuthRequest(t, router, "testuser", "user1")
	expectedToken, err := middleware.JWT{}.GenerateToken(&user)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestSendCoinHandlerIntegration(t *testing.T) {
	router, _, db := setupTestDB()
	defer clearDatabase(db)

	token := performAuthRequest(t, router, "user1", "user1")

	reqBody := map[string]interface{}{
		"receiver_username": "user2",
		"amount":            30.0,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBuffer(reqBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: token})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 30.0, resBody["money_amount"])
}

func TestBuyItemHandlerIntegration(t *testing.T) {
	router, _, db := setupTestDB()
	defer clearDatabase(db)

	token := performAuthRequest(t, router, "user1", "user1")

	req := httptest.NewRequest(http.MethodPost, "/api/buy/socks", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: token})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "socks", resBody["merch_name"])
}

func TestInfoHandlerIntegration(t *testing.T) {
	router, _, db := setupTestDB()
	defer clearDatabase(db)

	token := performAuthRequest(t, router, "user1", "user1")

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: token})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, resBody, "purchases")
	assert.Contains(t, resBody, "transactions")
}

func clearDatabase(db *gorm.DB) {
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM purchases")
	db.Exec("DELETE FROM merches")
	db.Exec("DELETE FROM users")
}
