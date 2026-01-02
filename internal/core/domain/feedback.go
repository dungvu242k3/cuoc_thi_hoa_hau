package domain

import "time"

type FeedbackType string

const (
	FeedbackProposal  FeedbackType = "PROPOSAL"  // Đề xuất
	FeedbackComplaint FeedbackType = "COMPLAINT" // Khiếu nại
	FeedbackRequest   FeedbackType = "REQUEST"   // Yêu cầu đặc biệt
)

type FeedbackStatus string

const (
	FeedbackPending  FeedbackStatus = "PENDING"  // Đang chờ xử lý
	FeedbackResolved FeedbackStatus = "RESOLVED" // Đã giải quyết
	FeedbackRejected FeedbackStatus = "REJECTED" // Từ chối
)

type Feedback struct {
	ID        string         `bson:"_id,omitempty" json:"id"`
	UserID    string         `bson:"user_id" json:"userId"` // ID của thí sinh gửi
	Type      FeedbackType   `bson:"type" json:"type"`
	Title     string         `bson:"title" json:"title"`
	Content   string         `bson:"content" json:"content"`
	Status    FeedbackStatus `bson:"status" json:"status"`
	CreatedAt time.Time      `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `bson:"updated_at" json:"updatedAt"`
}
