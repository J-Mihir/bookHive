package middleware

import (
	"context"
	"fmt" // Import fmt for logging
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

// We'll use an init function here as well to ensure both files see the same key.
func init() {
	secret := os.Getenv("JWT_SECRET_KEY")
	// --- DIAGNOSTIC LOG ---
	// This will print the key your middleware is using to the terminal on startup.
	fmt.Printf("INFO: Middleware JWT_SECRET_KEY being used: '%s'\n", secret)
	if secret == "" {
		fmt.Println("WARNING: Middleware JWT_SECRET_KEY not set. Using default.")
		secret = "default_insecure_secret_key"
	}
	jwtKey = []byte(secret)
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const userContextKey = contextKey("user")

// JWTMiddleware validates the token and adds user claims to the request context.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequired checks the role from the JWT claims in the context.
func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
			return
		}

		if claims.Role != "staff" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
