package middleware

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
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

func (JWT) ParseToken(tokenString string) (claims *domain.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
