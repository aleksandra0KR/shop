package domain

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	jwt.StandardClaims
}
