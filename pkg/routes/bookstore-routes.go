package routes

import (
	"github.com/J-Mihir/go-bookstore/pkg/controllers"
	"github.com/J-Mihir/go-bookstore/pkg/middleware"
	"github.com/gorilla/mux"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	// --- PUBLIC ROUTES ---
	// These routes do not require any authentication.
	router.HandleFunc("/books", controllers.GetBook).Methods("GET")
	router.HandleFunc("/books/{bookId}", controllers.GetBookById).Methods("GET")

	// --- PROTECTED ADMIN-ONLY ROUTES ---
	// Create a sub-router for routes that require a valid JWT and an admin role.
	adminRoutes := router.PathPrefix("/books").Subrouter()

	// Apply middlewares in the correct order.
	// 1. JWTMiddleware runs first to validate the token and add claims to the context.
	// 2. AdminRequired runs second to check the role from the claims in the context.
	adminRoutes.Use(middleware.JWTMiddleware)
	adminRoutes.Use(middleware.AdminRequired)

	// These handlers are now fully protected.
	adminRoutes.HandleFunc("", controllers.CreateBook).Methods("POST")
	adminRoutes.HandleFunc("/{bookId}", controllers.UpdateBook).Methods("PUT")
	adminRoutes.HandleFunc("/{bookId}", controllers.DeleteBook).Methods("DELETE")
}
