package routes

import (
	"log"

	"github.com/J-Mihir/go-bookstore/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterTransactionRoutes = func(router *mux.Router) {
	log.Println("Registering transaction routes...")
	router.HandleFunc("/transactions/borrow", controllers.BorrowBook).Methods("POST")
	router.HandleFunc("/transactions/{transactionId}/return", controllers.ReturnBook).Methods("PUT")
}
