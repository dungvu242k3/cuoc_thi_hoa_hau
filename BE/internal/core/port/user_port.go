package port

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
