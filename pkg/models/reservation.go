package models

import (
	"gorm.io/gorm"
)

// Reservation represents a user's request for a book that is currently borrowed.
type Reservation struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	User   User   `json:"user,omitempty"`
	BookID uint   `json:"book_id"`
	Book   Book   `json:"book,omitempty"`
	Status string `json:"status"` // e.g., "Pending", "Fulfilled", "Cancelled"
}
