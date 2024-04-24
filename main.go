package main

import (
	"encoding/json"
	"fmt"
	"log"
	"minna-style-hub/database"
	"minna-style-hub/handlers"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// Credentials struct for parsing login request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CustomClaims represents custom claims for JWT token
type CustomClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

// SecretKey for signing JWT tokens
var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))


// GenerateJWTToken generates a JWT token for the given username and isAdmin status
func GenerateJWTToken(username string, isAdmin bool) (string, error) {
	claims := CustomClaims{
		username,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LoginHandler handles user authentication and issues JWT token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		log.Fatal("ADMIN_USERNAME not found in .env file")
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD not found in .env file")
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(creds.Username)
	fmt.Println(creds.Password)
	fmt.Println(adminUsername)
	fmt.Println(adminPassword)

	// Perform authentication (replace this with your actual authentication logic)
	isAdmin := creds.Username == adminUsername && creds.Password == adminPassword

	// Generate JWT token
	if isAdmin {
		token, err := GenerateJWTToken(creds.Username, isAdmin)
		if err != nil {
			log.Println("Error generating JWT token:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Respond with JWT token
		response := map[string]interface{}{
			"token":   token,
			"isAdmin": isAdmin,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, "Invalid admin credentials", http.StatusBadRequest)
	}

}

// AuthMiddleware validates JWT token and authorizes users
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// Add custom logic to authorize users based on claims
			// For example, check if the user is an admin
			if !claims.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

func main() {
	r := mux.NewRouter()

	// Connect to MongoDB
	if err := database.ConnectToMongoDB(); err != nil {
		log.Fatal(err)
	}

	// Define routes
	r.HandleFunc("/items", handlers.GetAllItems).Methods("GET")
	r.HandleFunc("/item/{id}", handlers.GetItem).Methods("GET")
	r.HandleFunc("/feedback", handlers.GetFeedback).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.Handle("/items/add", AuthMiddleware(http.HandlerFunc(handlers.AddItem))).Methods("POST")
	r.Handle("/items/update", AuthMiddleware(http.HandlerFunc(handlers.UpdateItem))).Methods("PUT")
	r.Handle("/items/{id}", AuthMiddleware(http.HandlerFunc(handlers.DeleteItem))).Methods("DELETE")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
