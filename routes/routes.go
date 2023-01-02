package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	ProdService services.ProductServices
	Repo        repository.Repository
	ProdRepo    repository.RepositoryProduct
	Middleware  middleware.Middleware
}

func NewRoutes() Routes {
	return &App{}
}

var schema = `
CREATE TABLE IF NOT EXISTS product (
		id serial,
		product_name VARCHAR(20),
		product_description TEXT,
		price FLOAT,
		product_image VARCHAR(1000),
		created_at TIMESTAMP default CURRENT_TIMESTAMP,
		updated_at TIMESTAMP default CURRENT_TIMESTAMP
		)`

func (a *App) Run() {
	dsn := "user=postgres password=root dbname=db_golang sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Println(err.Error())
		return
	}
	db.MustExec(schema)
	a.Repo = repository.NewRepository(db)
	a.ProdRepo = repository.NewProductRepository(db)
	a.UserService = services.NewServices(a.Repo)
	a.ProdService = services.NewProductServices(a.ProdRepo)
	a.Middleware = middleware.NewMiddleware()
	a.SetupRouter()
	a.Router.Use(a.Middleware.LoggingMiddleware)
	log.Println("Starting Server")
	a.Logger.Fatal(http.ListenAndServe(":1323", a.Router))
}

func (a *App) SetupRouter() {
	a.Router = mux.NewRouter()
	var handlerfun handler.Handler = handler.NewHandler(a.UserService, a.Repo)
	var prodhandler handler.ProductHandler = handler.NewProductHandler(a.ProdService, a.ProdRepo)

	router := a.Router.PathPrefix("/users").Subrouter()
	router.Path("/getall").HandlerFunc(handlerfun.GetUsers).Methods("GET")
	router.Path("/create").HandlerFunc(handlerfun.CreateUser).Methods("POST")
	router.Path("/getbyId/").HandlerFunc(handlerfun.GetUserbyId).Methods("GET")
	router.Path("/update").HandlerFunc(handlerfun.UpdateUser).Methods("PUT")
	router.Path("/del").HandlerFunc(handlerfun.DeleteUser).Methods("DELETE")

	prod := a.Router.PathPrefix("/product").Subrouter()
	prod.Path("/getall").HandlerFunc(prodhandler.GetProducts).Methods("GET")
	prod.Path("/create").HandlerFunc(prodhandler.CreateProduct).Methods("POST")
	prod.Path("/getbyId/").HandlerFunc(prodhandler.GetProductbyId).Methods("GET")
	prod.Path("/update").HandlerFunc(prodhandler.UpdateProduct).Methods("PUT")
	prod.Path("/del").HandlerFunc(prodhandler.DeleteProduct).Methods("DELETE")
}
