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

var NewUser models.User

func GetUser(w http.ResponseWriter, r *http.Request) {
	newUsers := models.GetAllUsers()
	res, _ := json.Marshal(newUsers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	ID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userDetails, _ := models.GetUserById(ID)
	if userDetails.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res, _ := json.Marshal(userDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &models.User{}
	utils.ParseBody(r, newUser)

	// Call the updated CreateUser method which returns an error
	u, err := newUser.CreateUser()

	// If there was an error (e.g., duplicate email), send a 409 Conflict response
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(map[string]string{"error": "User with this email or membership ID already exists."})
		return
	}

	// If successful, send the created user
	res, _ := json.Marshal(u)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	ID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
	}
	user := models.DeleteUser(ID)
	res, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// First, get the user we want to update from the database
	vars := mux.Vars(r)
	userId := vars["userId"]
	ID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing userId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userDetails, db := models.GetUserById(ID)
	if userDetails.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Now, parse the incoming JSON into a temporary map
	// This allows us to check which fields were actually sent in the request
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check for each field in the map and update if it exists
	if name, ok := updateData["name"].(string); ok {
		userDetails.Name = name
	}
	if email, ok := updateData["email"].(string); ok {
		userDetails.Email = email
	}
	if membershipID, ok := updateData["membership_id"].(string); ok {
		userDetails.MembershipID = membershipID
	}
	if role, ok := updateData["role"].(string); ok {
		userDetails.Role = role
	}

	// This is the key fix: We check if "fines" exists in the request.
	// This works even if the value is 0.
	if fines, ok := updateData["fines"].(float64); ok {
		userDetails.Fines = fines
	}

	// Save the changes to the database
	db.Save(&userDetails)

	// Return the updated user details
	res, _ := json.Marshal(userDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
