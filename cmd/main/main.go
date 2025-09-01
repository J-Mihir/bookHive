package main

import (
	"log"
	"net/http"

	"github.com/J-Mihir/go-bookstore/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisterTransactionRoutes(r)
	routes.RegisterCategoryRoutes(r)
	routes.RegisterAuthRoutes(r)
	http.Handle("/", r)
	log.Println("Server running at http://localhost:9010")
	log.Fatal(http.ListenAndServe(":9010", r))
}
