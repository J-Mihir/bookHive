package models

import "gorm.io/gorm"

// Category represents a genre or classification for a book.
type Category struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique"`
	Books []Book `json:"-" gorm:"foreignKey:CategoryID"` // One-to-Many relationship, ignored in JSON response for simplicity
}

// CreateCategory creates a new category record in the database.
func (c *Category) CreateCategory() (*Category, error) {
	db := GetDB()
	result := db.Create(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}
