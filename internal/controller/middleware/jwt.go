package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"shop/domain"
	"time"
)

type JWT struct{}

func (JWT) GenerateToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &domain.Claims{
		Username: user.Username,
		UserGUID: user.Guid,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Guid,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Errorf("error generating token: %v", err)
	}

	return tokenString, err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("accessToken")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		username := claims["username"].(string)
		c.Set("username", username)
		c.Next()
	}
}
