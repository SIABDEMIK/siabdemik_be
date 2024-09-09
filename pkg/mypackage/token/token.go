package token

import (
	"Ayala-Crea/server-app-absensi/models"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key") // Ubah ini sesuai dengan kunci rahasia JWT Anda

// DecodeJWT adalah helper function untuk memvalidasi dan decode JWT token
func DecodeJWT(authHeader string) (*models.Claims, error) {
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	// Pastikan token menggunakan prefix "Bearer "
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return nil, errors.New("invalid authorization header format")
	}

	tokenString := strings.TrimSpace(parts[1])
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
