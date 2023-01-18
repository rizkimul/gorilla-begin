package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("secret")

type JWTClaim struct {
	jwt.RegisteredClaims
}
