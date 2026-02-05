package port

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, username, password, role string) (*domain.AuthClaims, string, error)
	Login(ctx context.Context, username, password string) (*domain.AuthClaims, string, error)
}
