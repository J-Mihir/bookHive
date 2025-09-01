package routes

import (
	"github.com/J-Mihir/go-bookstore/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users", controllers.GetUser).Methods("GET")
	router.HandleFunc("/users/{userId}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/users/{userId}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{userId}", controllers.DeleteUser).Methods("DELETE")
}
