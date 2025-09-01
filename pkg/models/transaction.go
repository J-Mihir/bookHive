package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID     uint       `json:"user_id"`
	User       User       `gorm:"foreignKey:UserID"`
	BookID     uint       `json:"book_id"`
	Book       Book       `gorm:"foreignKey:BookID"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"` // Pointer to handle null values
	Fine       float64    `json:"fine"`
}

// You can add functions here later for getting transaction history if needed.
// For now, the controller will handle creation and updates directly.
