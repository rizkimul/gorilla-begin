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
		PRIMARY KEY(id)
		);

CREATE TABLE IF NOT EXISTS cart (
		id serial,
		cart_name VARCHAR(50),
		PRIMARY KEY(id)
		);

CREATE TABLE IF NOT EXISTS shopping_cart (
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
	dsn := "user=postgres password=root dbname=db_golang sslmode=disable"
	Db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Println(err.Error())
		return
	}
	Db.MustExec(schema)
	a.Repo = repository.NewRepository(Db)
	a.ProdRepo = repository.NewProductRepository(Db)
	a.CartRepo = repository.NewCartRepository(Db)
	a.SPCartRepo = repository.NewSPCartRepository(Db)
	a.UserService = services.NewServices(a.Repo)
	a.ProdService = services.NewProductServices(a.ProdRepo)
	a.CartService = services.NewCartServices(a.CartRepo)
	a.SPCartService = services.NewSPCartServices(a.SPCartRepo)
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
	var carthandler handler.CartHandler = handler.NewCartHandler(a.CartService, a.CartRepo)
	var spcarthandler handler.SPCartHandler = handler.NewSPCartHandler(a.SPCartService, a.SPCartRepo)

	router := a.Router.PathPrefix("/users").Subrouter()
	router.Path("/getall").HandlerFunc(handlerfun.GetUsers).Methods("GET")
	router.Path("/create").HandlerFunc(handlerfun.CreateUser).Methods("POST")
	router.Path("/getbyId").HandlerFunc(handlerfun.GetUserbyId).Methods("GET")
	router.Path("/update").HandlerFunc(handlerfun.UpdateUser).Methods("PUT")
	router.Path("/del").HandlerFunc(handlerfun.DeleteUser).Methods("DELETE")

	prod := a.Router.PathPrefix("/product").Subrouter()
	prod.Path("/getall").HandlerFunc(prodhandler.GetProducts).Methods("GET")
	prod.Path("/create").HandlerFunc(prodhandler.CreateProduct).Methods("POST")
	prod.Path("/getbyId").HandlerFunc(prodhandler.GetProductbyId).Methods("GET")
	prod.Path("/update").HandlerFunc(prodhandler.UpdateProduct).Methods("PUT")
	prod.Path("/del").HandlerFunc(prodhandler.DeleteProduct).Methods("DELETE")

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
