package service

import (
	"context"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
)

type feedbackService struct {
	repo port.FeedbackRepository
}

func NewFeedbackService(repo port.FeedbackRepository) port.FeedbackService {
	return &feedbackService{repo: repo}
}

func (s *feedbackService) SendFeedback(ctx context.Context, userID string, f *domain.Feedback) error {
	// Assign UserID from context/auth
	f.UserID = userID

	// Default Status is PENDING
	f.Status = domain.FeedbackPending

	return s.repo.Create(ctx, f)
}

func (s *feedbackService) GetMyFeedbacks(ctx context.Context, userID string, limit, offset int64) ([]*domain.Feedback, int64, error) {
	// Basic validation
	const maxLimit = 50
	if limit <= 0 {
		limit = 10
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetListByUser(ctx, userID, limit, offset)
}
