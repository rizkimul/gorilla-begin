package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/rizkimul/gorilla-begin/v2/config"
	"github.com/rizkimul/gorilla-begin/v2/helper"
)

type Middleware interface {
	LoggingMiddleware(next http.Handler) http.Handler
	JWT(next http.Handler) http.Handler
}

type middleware struct {
	helper helper.Helper
}

func NewMiddleware(helper helper.Helper) Middleware {
	return &middleware{
		helper: helper,
	}
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func (m *middleware) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			res := map[string]interface{}{"message": "Invalid Token", "is_success": false, "status": "401"}
			m.helper.ResponseJSON(w, http.StatusUnauthorized, res)
			return
		}
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims := &config.JWTClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				res := map[string]interface{}{"message": "Signature Invalid", "is_success": false, "status": "401"}
				m.helper.ResponseJSON(w, http.StatusUnauthorized, res)
				return
			case jwt.ValidationErrorExpired:
				res := map[string]interface{}{"message": "Token expires", "is_success": false, "status": "401"}
				m.helper.ResponseJSON(w, http.StatusUnauthorized, res)
				return
			default:
				res := map[string]interface{}{"message": "Unauthorized", "is_success": false, "status": "401", "data": err.Error()}
				m.helper.ResponseJSON(w, http.StatusUnauthorized, res)
				return
			}
		}

		if !token.Valid {
			res := map[string]interface{}{"message": "Invalid Token", "is_success": false, "status": "401"}
			m.helper.ResponseJSON(w, http.StatusUnauthorized, res)
			return
		}

		ctx := context.WithValue(context.Background(), "UserInfo", claims)
		r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
