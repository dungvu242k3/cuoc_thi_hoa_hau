package port

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

// Repository Interface
type ContestantRepository interface {
	Create(ctx context.Context, c *domain.Contestant) error
	Update(ctx context.Context, c *domain.Contestant) error
	GetByUserID(ctx context.Context, userID string) (*domain.Contestant, error)
	GetByID(ctx context.Context, id string) (*domain.Contestant, error)
	// BChỉ trả về list đã lọc thông tin nhạy cảm
	GetPublicList(ctx context.Context, limit int64, offset int64) ([]*domain.Contestant, int64, error)
	GetPublicDetail(ctx context.Context, id string) (*domain.Contestant, error)
	CheckIdentifyCard(ctx context.Context, cardID string) (bool, error)
}

// Service Interface
type ContestantService interface {
	CreateProfile(ctx context.Context, userID string, input *domain.Contestant) (*domain.Contestant, error)
	UpdateProfile(ctx context.Context, userID string, input *domain.Contestant) (*domain.Contestant, error)
	SubmitProfile(ctx context.Context, userID string) error
	GetMyProfile(ctx context.Context, userID string) (*domain.Contestant, error)
	DeleteProfile(ctx context.Context, userID string) error
	GetPublicList(ctx context.Context, limit int64, offset int64) ([]*domain.Contestant, int64, error)
	GetPublicDetail(ctx context.Context, id string) (*domain.Contestant, error)
}
