package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rizkimul/gorilla-begin/v2/config"
)

type Token interface {
	CreateToken(ttl time.Duration) (string, error)
}

type token struct{}

func NewUtilsToken() Token {
	return &token{}
}

func (*token) CreateToken(ttl time.Duration) (string, error) {
	expTime := time.Now().UTC()
	claims := &config.JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime.Add(ttl)),
		},
	}
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		return "", err
	}

	return token, nil
}
