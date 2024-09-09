package data

import (
	"Ayala-Crea/server-app-absensi/models"
	"Ayala-Crea/server-app-absensi/pkg/mypackage/token"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func UploadExcel(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")

		// Gunakan helper DecodeJWT
		claims, err := token.DecodeJWT(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// id_penginputan diambil dari JWT
		idPenginputan := claims.IDPenginputan
		fmt.Printf("Validated id_penginputan: %d\n", idPenginputan)

		// AdminID diambil dari JWT
		adminID := claims.ID // Asumsikan ID di JWT adalah AdminID
		fmt.Printf("Validated AdminID: %d\n", adminID)

		// 3. Parse multipart form untuk mengunggah file
		err = r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// 4. Retrieve file dari form data
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 5. Buka file Excel menggunakan excelize
		f, err := excelize.OpenReader(file)
		if err != nil {
			http.Error(w, "Unable to read Excel file", http.StatusBadRequest)
			return
		}

		// 6. Dapatkan baris dari sheet pertama
		rows, err := f.GetRows(f.GetSheetName(0))
		if err != nil {
			http.Error(w, "Unable to read rows from Excel file", http.StatusBadRequest)
			return
		}

		// 7. Iterate melalui setiap baris dan simpan ke database
		for i, row := range rows {
			if i == 0 {
				// Lewatkan baris pertama jika itu adalah header
				continue
			}

			// Gunakan AdminID dari JWT
			student := models.StudentsEmployees{
				AdminID:     adminID, // Mengambil AdminID dari JWT token
				FullName:    row[1],
				Status:      row[2],
				Class:       row[3],
				NpkOrNpm:    row[4],
				PhoneNumber: row[5],
			}

			// Insert ke tabel students_employees
			query := `INSERT INTO students_employees (admin_id, full_name, status, class, npk_or_npm, phone_number)
                      VALUES ($1, $2, $3, $4, $5, $6)`

			_, err := db.Exec(query, student.AdminID, student.FullName, student.Status, student.Class, student.NpkOrNpm, student.PhoneNumber)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to insert record for row %d: %v", i+1, err), http.StatusInternalServerError)
				return
			}

			// Hash password dari NpkOrNpm
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.NpkOrNpm), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error hashing password", http.StatusInternalServerError)
				return
			}

			// Insert ke tabel users
			user := models.Users{
				IDRole:        2, // Misalnya IDRole adalah 1, sesuaikan dengan kebutuhan
				IDPenginputan: idPenginputan,
				Nama:          student.FullName,
				Username:      student.NpkOrNpm,       // Username tetap NpkOrNpm asli
				Password:      string(hashedPassword), // Password di-hash
				Email:         row[6],                 // Email diambil dari kolom lain (misal kolom ke-7)
				PhoneNumber:   student.PhoneNumber,
				CreatedAt:     time.Now(),
				IsActive:      true,
			}

			queryUser := `INSERT INTO users (id_role, id_penginputan, nama, username, password, email, phone_number, created_at, is_active)
                          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

			_, err = db.Exec(queryUser, user.IDRole, user.IDPenginputan, user.Nama, user.Username, user.Password, user.Email, user.PhoneNumber, user.CreatedAt, user.IsActive)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to insert user record for row %d: %v", i+1, err), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func convertToInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func GetAllStudentsEmployees(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")

		// Gunakan helper DecodeJWT
		_, err := token.DecodeJWT(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Query untuk mengambil semua data dari tabel students_employees
		query := `SELECT id, admin_id, full_name, status, class, npk_or_npm, phone_number FROM students_employees`

		// Eksekusi query
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Slice untuk menampung hasil query
		var studentsEmployees []models.StudentsEmployees

		// Iterasi melalui hasil query
		for rows.Next() {
			var student models.StudentsEmployees
			err := rows.Scan(&student.ID, &student.AdminID, &student.FullName, &student.Status, &student.Class, &student.NpkOrNpm, &student.PhoneNumber)
			if err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			studentsEmployees = append(studentsEmployees, student)
		}

		// Cek error setelah iterasi selesai
		if err = rows.Err(); err != nil {
			http.Error(w, "Error occurred during iteration", http.StatusInternalServerError)
			return
		}

		// Kembalikan hasil dalam bentuk JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(studentsEmployees)
	}
}

func GetDataByIdAdmin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")

		// Gunakan helper DecodeJWT
		claims, err := token.DecodeJWT(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// 3. Gunakan admin_id dari klaim JWT
		idAdmin := claims.ID

		// 4. Query untuk mendapatkan semua data berdasarkan admin_id
		query := `SELECT id, admin_id, full_name, status, class, npk_or_npm, phone_number FROM students_employees WHERE admin_id = $1`
		rows, err := db.Query(query, idAdmin)
		if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 5. Buat slice untuk menampung semua data
		var students []models.StudentsEmployees

		// 6. Iterate melalui hasil query dan tambahkan ke slice
		for rows.Next() {
			var student models.StudentsEmployees
			err := rows.Scan(&student.ID, &student.AdminID, &student.FullName, &student.Status, &student.Class, &student.NpkOrNpm, &student.PhoneNumber)
			if err != nil {
				http.Error(w, "Failed to scan record", http.StatusInternalServerError)
				return
			}
			students = append(students, student)
		}

		// 7. Periksa jika terjadi kesalahan saat iterasi
		if err = rows.Err(); err != nil {
			http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
			return
		}

		// 8. Kembalikan hasil sebagai JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(students)
	}
}

func CreateDataManual(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")

		// Gunakan helper DecodeJWT
		claims, err := token.DecodeJWT(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Ambil AdminID dari klaim JWT
		adminID := claims.ID

		// 3. Decode JSON dari body request
		var data models.StudentsEmployees
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Set nilai tambahan yang tidak diambil dari JSON
		data.AdminID = adminID
		// data.CreatedAt = time.Now()  // Set waktu sekarang
		// data.IsActive = true          // Default aktif

		// 5. Insert data ke dalam database
		query := `INSERT INTO students_employees (admin_id, full_name, status, class, npk_or_npm, phone_number)
		          VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = db.Exec(query, data.AdminID, data.FullName, data.Status, data.Class, data.NpkOrNpm, data.PhoneNumber)
		if err != nil {
			http.Error(w, "Failed to insert record: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 6. Set header response dan kirim response JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
