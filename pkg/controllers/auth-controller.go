package controllers

import (
	"encoding/json"
	"fmt" // Import the fmt package for printing
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/J-Mihir/go-bookstore/pkg/middleware" // Import middleware to use its Claims struct
	"github.com/J-Mihir/go-bookstore/pkg/models"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

// We'll use an init function to set the key and add our diagnostic print statement.
func init() {
	secret := os.Getenv("JWT_SECRET_KEY")
	// --- DIAGNOSTIC LOG ---
	// This will print the key your application is using to the terminal on startup.
	fmt.Printf("INFO: JWT_SECRET_KEY being used: '%s'\n", secret)
	if secret == "" {
		fmt.Println("WARNING: JWT_SECRET_KEY environment variable not set. Using a default, insecure key.")
		secret = "default_insecure_secret_key" // Fallback for safety, but shows an error state
	}
	jwtKey = []byte(secret)
}

// RegisterUser handles new user registration.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// The password hashing is handled by the BeforeSave hook in the User model.
	if _, err := user.CreateUser(); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusConflict)
		return
	}

	// Do not return the password in the response
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginUser handles user authentication and token generation.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	var user models.User
	if result := db.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the stored hashed password with the one provided in the request
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the JWT claims, using the struct from the middleware package
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with the claims and sign it with our secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Finally, send the token to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
