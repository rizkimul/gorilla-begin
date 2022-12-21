package main

import (
	"log"
	"net/http"

	"github.com/rizkimul/gorilla-begin/v2/middleware"
	"github.com/rizkimul/gorilla-begin/v2/routes"
)

func main() {
	a := routes.App{}
	a.SetupRouter()

	a.Router.Use(middleware.LoggingMiddleware)
	log.Println("Starting Server")
	a.Logger.Fatal(http.ListenAndServe(":1323", a.Router))
}
