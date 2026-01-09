package service

import (
	"context"
	"errors"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
)

type scoringService struct {
	repo           port.ScoreRepository
	contestantRepo port.ContestantRepository
}

func NewScoringService(repo port.ScoreRepository, contestantRepo port.ContestantRepository) port.ScoreService {
	return &scoringService{
		repo:           repo,
		contestantRepo: contestantRepo,
	}
}

func (s *scoringService) GetMyScores(ctx context.Context, userID string) ([]*domain.Score, error) {
	// 1. Get Contestant Identity
	contestant, err := s.contestantRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("không tìm thấy hồ sơ thí sinh")
	}

	// 2. Fetch Scores
	scores, err := s.repo.GetListByCandidateID(ctx, contestant.ID)
	if err != nil {
		return nil, err
	}

	return scores, nil
}
