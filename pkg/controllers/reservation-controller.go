package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/J-Mihir/go-bookstore/pkg/models"
)

// CreateReservation handles a user's request to reserve a book.
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	type ReservationRequest struct {
		UserID uint `json:"user_id"`
		BookID uint `json:"book_id"`
	}

	var req ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := models.GetDB()

	// 1. Validate the user
	user, _ := models.GetUserById(int64(req.UserID))
	if user.ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 2. Validate the book
	book, _ := models.GetBookById(int64(req.BookID))
	if book.ID == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// 3. Business Rule: A user can only reserve a book if it's currently "Borrowed".
	if book.Availability != "Borrowed" {
		http.Error(w, "This book is not currently borrowed and cannot be reserved.", http.StatusConflict)
		return
	}

	// 4. Check if the user already has a pending reservation for this book
	var existingReservation models.Reservation
	db.Where("user_id = ? AND book_id = ? AND status = ?", req.UserID, req.BookID, "Pending").First(&existingReservation)
	if existingReservation.ID != 0 {
		http.Error(w, "You already have a pending reservation for this book.", http.StatusConflict)
		return
	}

	// 5. Create the reservation
	reservation := models.Reservation{
		UserID: req.UserID,
		BookID: req.BookID,
		Status: "Pending",
	}

	if result := db.Create(&reservation); result.Error != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservation)
}
