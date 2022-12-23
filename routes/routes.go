package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rizkimul/gorilla-begin/v2/handler"
	"github.com/rizkimul/gorilla-begin/v2/middleware"
)

type App struct {
	Router *mux.Router
	Logger *log.Logger
}

func (a *App) Run() {
	a.SetupRouter()
	a.Router.Use(middleware.LoggingMiddleware)
	log.Println("Starting Server")
	a.Logger.Fatal(http.ListenAndServe(":1323", a.Router))
}

func (a *App) SetupRouter() {
	a.Router = mux.NewRouter()

	router := a.Router.PathPrefix("/users").Subrouter()
	router.Path("/getall").HandlerFunc(handler.GetUsers).Methods("GET")
	router.Path("/create").HandlerFunc(handler.CreateUser).Methods("POST")
	router.Path("/getbyId/").HandlerFunc(handler.GetUserbyId).Methods("GET")
	router.Path("/update").HandlerFunc(handler.UpdateUser).Methods("PUT")
	router.Path("/del").HandlerFunc(handler.DeleteUser).Methods("DELETE")
}
