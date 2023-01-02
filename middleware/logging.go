package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type Middleware interface {
	LoggingMiddleware(next http.Handler) http.Handler
}

type middleware struct{}

func NewMiddleware() Middleware {
	return &middleware{}
}

func (*middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
