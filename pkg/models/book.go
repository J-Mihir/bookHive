package models

import (
	"github.com/J-Mihir/go-bookstore/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name         string `json:"name"`
	Author       string `json:"author"`
	Publication  string `json:"publication"`
	ISBN         string `json:"isbn" gorm:"unique"` // unique for each edition
	Genre        string `json:"genre"`
	Edition      string `json:"edition"`
	Copies       int    `json:"copies"` // total copies in library
	Availability string `json:"availability"`
	CategoryID   uint   `json:"category_id"` // Available, Borrowed, Reserved
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{}, &Transaction{}, &Category{}, &Book{})
}

func (b *Book) CreateBook() (*Book, error) {
	result := GetDB().Create(&b)
	if result.Error != nil {
		return nil, result.Error
	}
	return b, nil
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	result := db.Where("id = ?", Id).First(&getBook)
	return &getBook, result
}

func DeleteBook(Id int64) Book {
	var book Book
	db.Where("id = ?", Id).Delete(&book)
	return book
}
