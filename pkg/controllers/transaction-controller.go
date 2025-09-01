package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/J-Mihir/go-bookstore/pkg/models"
	"github.com/gorilla/mux"
)

const maxBorrowLimit = 5 // Business Rule: A user can borrow a maximum of 5 books

// BorrowBook handles the logic for a user borrowing a book.
func BorrowBook(w http.ResponseWriter, r *http.Request) {
	// This function remains unchanged
	type BorrowRequest struct {
		UserID uint `json:"user_id"`
		BookID uint `json:"book_id"`
	}

	var req BorrowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := models.GetDB()

	user, _ := models.GetUserById(int64(req.UserID))
	if user.ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	book, _ := models.GetBookById(int64(req.BookID))
	if book.ID == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	if book.Availability != "Available" {
		http.Error(w, "Book is currently not available", http.StatusConflict)
		return
	}

	var currentBorrows int64
	db.Model(&models.Transaction{}).Where("user_id = ? AND return_date IS NULL", req.UserID).Count(&currentBorrows)

	if currentBorrows >= maxBorrowLimit {
		errorMsg := fmt.Sprintf("Borrow limit of %d books reached", maxBorrowLimit)
		http.Error(w, errorMsg, http.StatusForbidden)
		return
	}

	transaction := models.Transaction{
		UserID:     req.UserID,
		BookID:     req.BookID,
		BorrowDate: time.Now(),
		DueDate:    time.Now().AddDate(0, 0, 14),
	}
	if result := db.Create(&transaction); result.Error != nil {
		http.Error(w, "Failed to create transaction: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	book.Availability = "Borrowed"
	db.Save(&book)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

// ReturnBook is updated to handle reservations.
func ReturnBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionIdStr := vars["transactionId"]
	transactionID, err := strconv.ParseUint(transactionIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	var transaction models.Transaction
	if result := db.First(&transaction, transactionID); result.Error != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if transaction.ReturnDate != nil {
		http.Error(w, "Book has already been returned", http.StatusConflict)
		return
	}

	returnDate := time.Now()
	transaction.ReturnDate = &returnDate

	if returnDate.After(transaction.DueDate) {
		daysOverdue := int(returnDate.Sub(transaction.DueDate).Hours() / 24)
		transaction.Fine = float64(daysOverdue * 1)
	}
	db.Save(&transaction)

	// --- NEW LOGIC: Check for reservations when a book is returned ---
	var reservation models.Reservation
	// Find the oldest pending reservation for this book
	result := db.Where("book_id = ? AND status = ?", transaction.BookID, "Pending").Order("created_at asc").First(&reservation)

	book, _ := models.GetBookById(int64(transaction.BookID))

	if result.Error == nil {
		// A reservation was found, so fulfill it.
		reservation.Status = "Fulfilled"
		db.Save(&reservation)
		// Set the book's status to "Reserved" for the next user.
		book.Availability = "Reserved"
	} else {
		// No reservations found, make the book generally available.
		book.Availability = "Available"
	}
	db.Save(&book)
	// --- END OF NEW LOGIC ---

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
