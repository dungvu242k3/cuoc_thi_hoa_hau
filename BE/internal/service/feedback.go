/*
Package service đảm nhận Business Logic.
Tác dụng của feedback.go:
- Gói gọn logic gửi và xem danh sách Phản hồi/Góp ý của người dùng/thí sinh.
*/
package service

import (
	"context"

	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"
)

type feedbackService struct {
	repo types.FeedbackRepository
}

func NewFeedbackService(repo types.FeedbackRepository) types.FeedbackService {
	return &feedbackService{repo: repo}
}

func (s *feedbackService) SendFeedback(ctx context.Context, userID string, f *model.Feedback) error {
	f.ContestantID = userID
	f.Status = model.FeedbackStatusPending
	return s.repo.Create(ctx, f)
}

// GetMyFeedbacks chặn max limit là 50 để chống việc cào data (scraping) làm tắc nghẽn Database.
func (s *feedbackService) GetMyFeedbacks(ctx context.Context, userID string, limit, offset int64) ([]*model.Feedback, int64, error) {
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
