package routes

import (
	"github.com/J-Mihir/go-bookstore/pkg/controllers"
	"github.com/J-Mihir/go-bookstore/pkg/middleware"
	"github.com/gorilla/mux"
)

var RegisterReservationRoutes = func(router *mux.Router) {
	// We'll protect reservation creation to ensure a user is identified.
	reservationRoutes := router.PathPrefix("/reservations").Subrouter()
	// Reusing middleware. A user must be identified via X-User-ID to make a reservation.
	reservationRoutes.Use(middleware.AdminRequired)

	reservationRoutes.HandleFunc("", controllers.CreateReservation).Methods("POST")
}
