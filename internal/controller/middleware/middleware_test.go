package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"shop/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Setup() *gin.Engine {
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	return router
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: ""})
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: "InvalidToken"})
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	router := Setup()

	user := domain.User{Username: "test", Password: "test"}
	validToken, err := JWT{}.GenerateToken(&user)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: validToken})
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
