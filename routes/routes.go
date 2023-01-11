package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rizkimul/gorilla-begin/v2/handler"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/middleware"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/services"
)

type Routes interface {
	Run()
}

type App struct {
	Router        *mux.Router
	Logger        *log.Logger
	UserService   services.Services
	ProdService   services.ProductServices
	CartService   services.CartServices
	SPCartService services.SPCartServices
	Repo          repository.Repository
	ProdRepo      repository.RepositoryProduct
	CartRepo      repository.CartRepository
	SPCartRepo    repository.SPCartRepository
	Middleware    middleware.Middleware
	Helper        helper.Helper
}

func NewRoutes() Routes {
	return &App{}
}

var (
	_       = godotenv.Load(".env")
	dbUser  = os.Getenv("DATABASE_USER")
	dbPass  = os.Getenv("DATABASE_PASSWORD")
	dbName  = os.Getenv("DATABASE_NAME")
	sslMode = os.Getenv("SSL_MODE")
)

var schema = `
CREATE TABLE IF NOT EXISTS product (
		id serial,
		product_name VARCHAR(20),
		product_description TEXT,
		price FLOAT,
		product_image VARCHAR(1000),
		created_at TIMESTAMP default CURRENT_TIMESTAMP,
		PRIMARY KEY(id)
		);

CREATE TABLE IF NOT EXISTS cart (
		id serial,
		cart_name VARCHAR(50),
		PRIMARY KEY(id)
		);

CREATE TABLE IF NOT EXISTS shopping_cart (
		id serial,
		cart_id INT,
		product_id INT,
		qty_product INT,
		total_price FLOAT,
		CONSTRAINT fk_cart
			FOREIGN KEY(cart_id)
				REFERENCES cart(id),
		CONSTRAINT fk_product
			FOREIGN KEY(product_id)
				REFERENCES product(id)
		)`

func (a *App) Run() {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPass, dbName, sslMode)
	Db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Println(err.Error())
		return
	}
	Db.MustExec(schema)
	a.Repo = repository.NewRepository(Db)
	a.ProdRepo = repository.NewProductRepository(Db)
	a.SPCartRepo = repository.NewSPCartRepository(Db, a.ProdRepo)
	a.CartRepo = repository.NewCartRepository(Db, a.SPCartRepo, a.ProdRepo)
	a.UserService = services.NewServices(a.Repo)
	a.ProdService = services.NewProductServices(a.ProdRepo)
	a.CartService = services.NewCartServices(a.CartRepo, a.SPCartRepo)
	a.SPCartService = services.NewSPCartServices(a.SPCartRepo)
	a.Helper = helper.NewHelper()
	a.Middleware = middleware.NewMiddleware(a.Helper)
	a.SetupRouter()
	a.Router.Use(a.Middleware.LoggingMiddleware)
	log.Println("Starting Server")
	a.Logger.Fatal(http.ListenAndServe(":1323", a.Router))
}

func (a *App) SetupRouter() {
	a.Router = mux.NewRouter()
	var handlerfun handler.Handler = handler.NewHandler(a.UserService, a.Repo, a.Helper)
	var prodhandler handler.ProductHandler = handler.NewProductHandler(a.ProdService, a.ProdRepo, a.Helper)
	var carthandler handler.CartHandler = handler.NewCartHandler(a.CartService, a.CartRepo, a.Helper)
	var spcarthandler handler.SPCartHandler = handler.NewSPCartHandler(a.SPCartService, a.SPCartRepo, a.Helper)

	router := a.Router.PathPrefix("/users").Subrouter()
	router.Path("/getall").HandlerFunc(handlerfun.GetUsers).Methods("GET")
	router.Path("/register").HandlerFunc(handlerfun.CreateUser).Methods("POST")
	router.Path("/login").HandlerFunc(handlerfun.Login).Methods("POST")
	router.Path("/getbyId").HandlerFunc(handlerfun.GetUserbyId).Methods("GET")
	router.Path("/update").HandlerFunc(handlerfun.UpdateUser).Methods("PUT")
	router.Path("/del").HandlerFunc(handlerfun.DeleteUser).Methods("DELETE")
	router.Path("/print").HandlerFunc(handlerfun.PrinData).Methods("GET")

	prod := a.Router.PathPrefix("/product").Subrouter()
	prod.Path("/getall").HandlerFunc(prodhandler.GetProducts).Methods("GET")
	prod.Path("/create").HandlerFunc(prodhandler.CreateProduct).Methods("POST")
	prod.Path("/getbyId").HandlerFunc(prodhandler.GetProductbyId).Methods("GET")
	prod.Path("/update").HandlerFunc(prodhandler.UpdateProduct).Methods("PUT")
	prod.Path("/del").HandlerFunc(prodhandler.DeleteProduct).Methods("DELETE")
	prod.Use(a.Middleware.JWT)

	cart := a.Router.PathPrefix("/cart").Subrouter()
	cart.Path("/getall").HandlerFunc(carthandler.GetCarts).Methods("GET")
	cart.Path("/create").HandlerFunc(carthandler.CreateCart).Methods("POST")
	cart.Path("/getbyId").HandlerFunc(carthandler.GetCartbyId).Methods("GET")
	cart.Path("/update").HandlerFunc(carthandler.UpdateCart).Methods("PUT")
	cart.Path("/del").HandlerFunc(carthandler.DeleteCart).Methods("DELETE")

	spcart := a.Router.PathPrefix("/shoppingcart").Subrouter()
	spcart.Path("/getall").HandlerFunc(spcarthandler.GetSPCarts).Methods("GET")
	spcart.Path("/create").HandlerFunc(spcarthandler.CreateSPCart).Methods("POST")
	spcart.Path("/getbyId").HandlerFunc(spcarthandler.GetSPCartbyId).Methods("GET")
	spcart.Path("/update").HandlerFunc(spcarthandler.UpdateSPCart).Methods("PUT")
	spcart.Path("/del").HandlerFunc(spcarthandler.DeleteSPCart).Methods("DELETE")
}
