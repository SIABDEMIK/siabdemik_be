package handlers

import (
	"Ayala-Crea/server-app-absensi/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

func RegisterUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.Users

		// Decode JSON request body into user struct
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Debugging: Log values to ensure they are parsed correctly
		log.Printf("Received IDRole: %d, IDPenginputan: %d\n", user.IDRole, user.IDPenginputan)

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		// Insert user into the database
		query := `INSERT INTO users (id_role, id_penginputan, nama, username, password, email, phone_number)
                  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, is_active`
		err = db.QueryRow(query, user.IDRole, user.IDPenginputan, user.Nama, user.Username, user.Password, user.Email, user.PhoneNumber).Scan(&user.ID, &user.CreatedAt, &user.IsActive)
		if err != nil {
			log.Printf("Error inserting user: %v\n", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Return the created user as a JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.Users
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Decode the JSON request
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Fetch the user from the database
		query := "SELECT id, id_role, id_penginputan, password, nama FROM users WHERE username=$1"
		row := db.QueryRow(query, input.Username)
		err = row.Scan(&user.ID, &user.IDRole, &user.IDPenginputan, &user.Password, &user.Nama)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Compare the provided password with the hashed password in the database
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Generate JWT Token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &models.Claims{
			ID:            user.ID,
			IDRole:        user.IDRole,
			IDPenginputan: user.IDPenginputan,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Create the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Return the token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token": tokenString,
			"nama":  user.Nama,
		})
	}
}
