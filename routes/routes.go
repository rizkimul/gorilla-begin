package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rizkimul/gorilla-begin/v2/handler"
)

type App struct {
	Router *mux.Router
	Logger *log.Logger
}

func (a *App) SetupRouter() {
	a.Router = mux.NewRouter()

	router := a.Router.PathPrefix("/users").Subrouter()
	router.Path("").HandlerFunc(handler.GetUsers).Methods("GET")
	router.Path("").HandlerFunc(handler.CreateUser).Methods("POST")
	router.Path("/get").HandlerFunc(handler.GetUserbyId).Methods("GET")
	router.Path("/update").HandlerFunc(handler.UpdateUser).Methods("PUT")
	router.Path("/del").HandlerFunc(handler.DeleteUser).Methods("DELETE")
}

func (a *App) CreateLoggingRouter(out io.Writer) http.Handler {
	return handlers.LoggingHandler(out, a.Router)
}
