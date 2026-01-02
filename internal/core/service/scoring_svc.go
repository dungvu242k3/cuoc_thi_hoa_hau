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

func NewScoringService(repo port.ScoreRepository, contestantRepo port.ContestantRepository) port.ScoringService {
	return &scoringService{
		repo:           repo,
		contestantRepo: contestantRepo,
	}
}

func (s *scoringService) GetMyScore(ctx context.Context, userID string) (*domain.Score, error) {
	// 1. Get Contestant Identity
	contestant, err := s.contestantRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("không tìm thấy hồ sơ thí sinh")
	}

	// 2. Fetch Score
	score, err := s.repo.GetByCandidateID(ctx, contestant.ID)
	if err != nil {
		return nil, err
	}

	if score == nil {
		return nil, errors.New("chưa có điểm số") // Or return empty score depending on requirements
	}

	// 3. Security: No masking needed for Candidate themselves as per requirements,
	// but if there were Examiner comments, we would remove them here.
	// For now, we return the full score as it contains only numbers and criteria.

	return score, nil
}
