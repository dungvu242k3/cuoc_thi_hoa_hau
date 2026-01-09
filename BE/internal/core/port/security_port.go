package port

import (
	"cuoc_thi_hoa_hau/internal/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PasswordEncoder interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type TokenProvider interface {
	Generate(user *domain.User, duration time.Duration) (*domain.AuthClaims, string, error)
	Validate(tokenString string) (jwt.MapClaims, error)
}
