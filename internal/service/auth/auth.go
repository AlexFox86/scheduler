package auth

import (
	"fmt"
	"os"

	"github.com/AlexFox86/scheduler/internal/pkg/token"
)

// Service provides authentication and login methods.
type Service struct {
	jwtSecret []byte
}

// New creates a new authentication service
func New(jwtSecret string) *Service {
	return &Service{
		jwtSecret: []byte(jwtSecret),
	}
}

// Password returns the 'password' field
func (s *Service) Password() string {
	return os.Getenv("TODO_PASSWORD")
}

// Login performs user authentication
func (s *Service) Login(password string) (string, error) {
	if password != s.Password() {
		return "", fmt.Errorf("Invalid password")
	}

	token, err := s.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateToken creates a JWT token
func (s *Service) GenerateToken() (string, error) {
	return token.Generate(s.jwtSecret)
}

// ValidateToken checks the JWT token
func (s *Service) ValidateToken(reqToken string) error {
	return token.Validate(reqToken, s.jwtSecret)
}
