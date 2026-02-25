/*
Package model định nghĩa cấu trúc Hòm thư góp ý/Khiếu nại.
Tác dụng:
- Lưu trữ các ý kiến của thí sinh/người dùng gửi về Ban Tổ Chức (được lưu vào collection feedbacks).
*/
package model

import "time"

type FeedbackType string

const (
	FeedbackTypeComplaint FeedbackType = "complaint" // Khiếu nại
	FeedbackTypeProposal  FeedbackType = "proposal"  // Đề xuất
	FeedbackTypeSupport   FeedbackType = "support"   // Hỗ trợ
)

type FeedbackStatus string

const (
	FeedbackStatusPending  FeedbackStatus = "pending"
	FeedbackStatusResolved FeedbackStatus = "resolved"
	FeedbackStatusClosed   FeedbackStatus = "closed"
)

// Feedback là bản ghi phản hồi của người dùng.
type Feedback struct {
	ID           string         `bson:"_id,omitempty" json:"id"`
	ContestantID string         `bson:"contestant_id" json:"contestantId"`
	Title        string         `bson:"title" json:"title"`
	Content      string         `bson:"content" json:"content"`
	Type         FeedbackType   `bson:"type" json:"type"`
	Status       FeedbackStatus `bson:"status" json:"status"`         // Trạng thái: đang chờ (pending), đã giải quyết (resolved)
	Reply        string         `bson:"reply,omitempty" json:"reply"` // Lời đáp trả lại từ Quản trị viên

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
