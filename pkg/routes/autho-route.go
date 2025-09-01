package routes

import (
	"github.com/J-Mihir/go-bookstore/pkg/controllers"
	"github.com/gorilla/mux"
)

// RegisterAuthRoutes sets up the public routes for user registration and login.
var RegisterAuthRoutes = func(router *mux.Router) {
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
}
