package api

import (
	"database/sql"

	handlers "Ayala-Crea/server-app-absensi/api/handlers"
	"Ayala-Crea/server-app-absensi/api/handlers/data"
	"Ayala-Crea/server-app-absensi/api/handlers/role"

	"github.com/gorilla/mux"
)

func AllRoutes(r *mux.Router, db *sql.DB) {
	// Register the POST route for user registration
	r.HandleFunc("/register", handlers.RegisterUser(db)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(db)).Methods("POST")

	// data karyawan and mahasiswa
	r.HandleFunc("/upload", data.UploadExcel(db)).Methods("POST")
	r.HandleFunc("/data", data.GetAllStudentsEmployees(db)).Methods("GET")
	r.HandleFunc("/data/mahasiswa", data.GetDataByIdAdmin(db)).Methods("GET")
	r.HandleFunc("/data/input", data.CreateDataManual(db)).Methods("POST")

	// role
	r.HandleFunc("/insert/role", role.InsertRole(db)).Methods("POST")
}
