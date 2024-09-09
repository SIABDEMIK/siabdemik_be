package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	App struct {
		Name        string `mapstructure:"name"`
		Environment string `mapstructure:"environment"`
		Port        int    `mapstructure:"port"`
	} `mapstructure:"app"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SslMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
}

type Users struct {
	ID            int       `json:"id"`
	IDRole        int       `json:"id_role"`        // Perbarui tag JSON
	IDPenginputan int       `json:"id_penginputan"` // Perbarui tag JSON
	Nama          string    `json:"nama"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	IsActive      bool      `json:"is_active"`
}

type Claims struct {
	ID            int `json:"id"`
	IDRole        int `json:"id_role"`
	IDPenginputan int `json:"id_penginputan"`
	jwt.StandardClaims
}

type StudentsEmployees struct {
	ID          int       `json:"id"`
	AdminID     int       `json:"admin_id"`
	FullName    string    `json:"full_name"`
	Status      string    `json:"status"`
	Class       string    `json:"class"`
	NpkOrNpm    string    `json:"npk_or_npm"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
}
