package main

import (
	"github.com/rizkimul/gorilla-begin/v2/routes"
)

var route routes.Routes = routes.NewRoutes()

func main() {
	route.Run()
}
