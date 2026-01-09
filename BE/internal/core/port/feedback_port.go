package port

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

// FeedbackRepository defines storage operations for Feedbacks.
type FeedbackRepository interface {
	// Create persists a new feedback into storage.
	Create(ctx context.Context, f *domain.Feedback) error
	// GetListByUser retrieves feedbacks for a specific user.
	GetListByUser(ctx context.Context, userID string, limit int64, offset int64) ([]*domain.Feedback, int64, error)
}

// FeedbackService defines business logic for Feedback management.
type FeedbackService interface {
	// SendFeedback allows a candidate to submit a feedback.
	SendFeedback(ctx context.Context, userID string, f *domain.Feedback) error
	// GetMyFeedbacks retrieves feedbacks submitted by the requesting user.
	GetMyFeedbacks(ctx context.Context, userID string, limit int64, offset int64) ([]*domain.Feedback, int64, error)
}
