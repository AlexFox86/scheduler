package token

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// Generate creates a JWT token
func Generate(jwtSecret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token.SignedString(jwtSecret)
}

// Validate checks the JWT token
func Validate(tokenString string, jwtSecret []byte) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("parse token: %w", err)
	}

	return nil
}
