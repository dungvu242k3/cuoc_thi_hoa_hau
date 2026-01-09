package domain

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

type Feedback struct {
	ID           string         `bson:"_id,omitempty" json:"id"`
	ContestantID string         `bson:"contestant_id" json:"contestantId"`
	Title        string         `bson:"title" json:"title"`
	Content      string         `bson:"content" json:"content"`
	Type         FeedbackType   `bson:"type" json:"type"`
	Status       FeedbackStatus `bson:"status" json:"status"`
	Reply        string         `bson:"reply,omitempty" json:"reply"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
