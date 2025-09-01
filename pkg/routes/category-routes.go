package routes

import (
	"github.com/J-Mihir/go-bookstore/pkg/controllers"
	"github.com/J-Mihir/go-bookstore/pkg/middleware"
	"github.com/gorilla/mux"
)

var RegisterCategoryRoutes = func(router *mux.Router) {
	// --- PUBLIC ROUTES for categories ---
	router.HandleFunc("/categories", controllers.GetAllCategories).Methods("GET")
	router.HandleFunc("/categories/{categoryId}", controllers.GetCategoryById).Methods("GET")

	// --- ADMIN-ONLY ROUTES for categories ---
	adminRoutes := router.PathPrefix("/categories").Subrouter()
	adminRoutes.Use(middleware.AdminRequired)

	adminRoutes.HandleFunc("", controllers.CreateCategory).Methods("POST")
	adminRoutes.HandleFunc("/{categoryId}", controllers.UpdateCategory).Methods("PUT")
	adminRoutes.HandleFunc("/{categoryId}", controllers.DeleteCategory).Methods("DELETE")
}
