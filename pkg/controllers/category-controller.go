package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/J-Mihir/go-bookstore/pkg/models"
	"github.com/J-Mihir/go-bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

// CreateCategory handles the creation of a new book category.
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	newCategory := &models.Category{}
	utils.ParseBody(r, newCategory)

	c, err := newCategory.CreateCategory()
	if err != nil {
		http.Error(w, "Failed to create category: "+err.Error(), http.StatusConflict)
		return
	}

	res, _ := json.Marshal(c)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetAllCategories retrieves all categories.
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	db := models.GetDB()
	var categories []models.Category
	db.Find(&categories)

	res, _ := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetCategoryById retrieves a single category by its ID.
func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryIdStr := vars["categoryId"]
	ID, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	var category models.Category
	if result := db.First(&category, ID); result.Error != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	res, _ := json.Marshal(category)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UpdateCategory updates a category's name.
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryIdStr := vars["categoryId"]
	ID, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	updateData := &models.Category{}
	utils.ParseBody(r, updateData)

	db := models.GetDB()
	var category models.Category
	if result := db.First(&category, ID); result.Error != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category.Name = updateData.Name
	db.Save(&category)

	res, _ := json.Marshal(category)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteCategory removes a category.
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryIdStr := vars["categoryId"]
	ID, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	db := models.GetDB()

	// Safety check: Prevent deleting a category if it still has books.
	var bookCount int64
	db.Model(&models.Book{}).Where("category_id = ?", ID).Count(&bookCount)
	if bookCount > 0 {
		http.Error(w, "Cannot delete category: it is still associated with books", http.StatusConflict)
		return
	}

	var category models.Category
	if result := db.First(&category, ID); result.Error != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	db.Delete(&category)

	w.WriteHeader(http.StatusNoContent) // 204 No Content is a good response for a successful delete
}
