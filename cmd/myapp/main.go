package main

import (
	"Ayala-Crea/server-app-absensi/api"
	"Ayala-Crea/server-app-absensi/pkg/config"
	"Ayala-Crea/server-app-absensi/pkg/cors"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name, cfg.Database.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping the database: %v", err)
	}

	// Create a new router
	r := mux.NewRouter()

	// Register API routes with the database connection
	api.AllRoutes(r, db)

	// Apply the CORS middleware to the router
	corsRouter := cors.CORSMiddleware(r)

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("Server is running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, corsRouter))
}
