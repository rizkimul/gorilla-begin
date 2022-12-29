package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rizkimul/gorilla-begin/v2/handler"
	"github.com/rizkimul/gorilla-begin/v2/middleware"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/services"
)

type Routes interface {
	Run()
}

type App struct {
	Router      *mux.Router
	Logger      *log.Logger
	UserService services.Services
	Repo        repository.Repository
}

func NewRoutes() Routes {
	return &App{}
}

func (a *App) Run() {
	a.Repo = repository.NewRepository()
	a.UserService = services.NewServices(a.Repo)
	a.SetupRouter()
	a.Router.Use(middleware.LoggingMiddleware)
	log.Println("Starting Server")
	a.Logger.Fatal(http.ListenAndServe(":1323", a.Router))
}

func (a *App) SetupRouter() {
	a.Router = mux.NewRouter()
	var handlerfun handler.Handler = handler.NewHandler(a.UserService, a.Repo)

	router := a.Router.PathPrefix("/users").Subrouter()
	router.Path("/getall").HandlerFunc(handlerfun.GetUsers).Methods("GET")
	router.Path("/create").HandlerFunc(handlerfun.CreateUser).Methods("POST")
	router.Path("/getbyId/").HandlerFunc(handlerfun.GetUserbyId).Methods("GET")
	router.Path("/update").HandlerFunc(handlerfun.UpdateUser).Methods("PUT")
	router.Path("/del").HandlerFunc(handlerfun.DeleteUser).Methods("DELETE")
}
