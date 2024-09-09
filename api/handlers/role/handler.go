package role

import (
	"Ayala-Crea/server-app-absensi/models"
	"Ayala-Crea/server-app-absensi/pkg/mypackage/token"
	"database/sql"
	"encoding/json"
	"net/http"
)

var jwtKey = []byte("your_secret_key")

func InsertRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Gunakan helper DecodeJWT
		claims, err := token.DecodeJWT(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Cek apakah role adalah 3
		if claims.IDRole != 3 {
			http.Error(w, "You do not have permission to insert data", http.StatusForbidden)
			return
		}

		var role models.Role
		err = json.NewDecoder(r.Body).Decode(&role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSER INTO role (nama, description) VALUES ($1, $2)`
		_, err = db.Exec(query, role.Name, role.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
