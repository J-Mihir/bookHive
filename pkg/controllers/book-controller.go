package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/J-Mihir/go-bookstore/pkg/models"
	"github.com/J-Mihir/go-bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

// GetBook retrieves all books
func GetBook(w http.ResponseWriter, r *http.Request) {
	newBooks := models.GetAllBooks()
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetBookById retrieves a single book by its ID
func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
	}
	bookDetails, _ := models.GetBookById(ID)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// CreateBook adds a new book to the database
func CreateBook(w http.ResponseWriter, r *http.Request) {
	createBook := &models.Book{}
	utils.ParseBody(r, createBook)

	// Check for required fields to prevent creating an empty book.
	if createBook.Name == "" || createBook.ISBN == "" || createBook.CategoryID == 0 {
		http.Error(w, "Missing required fields: name, isbn, and category_id are required", http.StatusBadRequest)
		return
	}

	// Set availability based on the number of copies
	if createBook.Copies > 0 {
		createBook.Availability = "Available"
	} else {
		createBook.Availability = "Not Available"
	}

	// Validate the Category ID to ensure it exists
	var category models.Category
	db := models.GetDB()
	result := db.First(&category, createBook.CategoryID)
	if result.Error != nil || result.RowsAffected == 0 {
		http.Error(w, "Invalid Category ID provided", http.StatusBadRequest)
		return
	}

	b, err := createBook.CreateBook()
	if err != nil {
		http.Error(w, "Failed to create book: "+err.Error(), http.StatusConflict)
		return
	}

	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteBook removes a book from the database
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
	}
	book := models.DeleteBook(ID)
	res, _ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UpdateBook modifies an existing book's details
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	updateBook := &models.Book{}
	utils.ParseBody(r, updateBook)

	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	bookDetails, db := models.GetBookById(ID)
	if bookDetails.ID == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Update fields only if they are provided in the request
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	if updateBook.ISBN != "" {
		bookDetails.ISBN = updateBook.ISBN
	}
	if updateBook.Genre != "" {
		bookDetails.Genre = updateBook.Genre
	}
	if updateBook.Edition != "" {
		bookDetails.Edition = updateBook.Edition
	}
	if updateBook.Copies > 0 { // Check greater than 0 to allow setting copies
		bookDetails.Copies = updateBook.Copies
	}
	if updateBook.Availability != "" {
		bookDetails.Availability = updateBook.Availability
	}
	if updateBook.CategoryID != 0 {
		bookDetails.CategoryID = updateBook.CategoryID
	}

	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
