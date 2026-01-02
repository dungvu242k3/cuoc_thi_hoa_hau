package port

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

type ScoreRepository interface {
	GetByCandidateID(ctx context.Context, candidateID string) (*domain.Score, error)
	Upsert(ctx context.Context, score *domain.Score) error
}

type ScoringService interface {
	GetMyScore(ctx context.Context, userID string) (*domain.Score, error)
}
