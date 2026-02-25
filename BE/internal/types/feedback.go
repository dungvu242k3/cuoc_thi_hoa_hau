/*
Package types định nghĩa Interface cho hòm thư Góp ý.
Tác dụng:
- FeedbackRepository: Nơi lưu trữ ý kiến vào DB.
- FeedbackService: Nơi kiểm duyệt ý kiến, gửi email thông báo (nếu có).
*/
package types

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
)

type FeedbackRepository interface {
	Create(ctx context.Context, f *model.Feedback) error
	GetListByUser(ctx context.Context, userID string, limit int64, offset int64) ([]*model.Feedback, int64, error)
}

type FeedbackService interface {
	SendFeedback(ctx context.Context, userID string, f *model.Feedback) error
	GetMyFeedbacks(ctx context.Context, userID string, limit int64, offset int64) ([]*model.Feedback, int64, error)
}
