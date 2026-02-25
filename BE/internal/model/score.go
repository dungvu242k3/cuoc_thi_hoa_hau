/*
Package model định nghĩa cấu trúc Phiếu Điểm.
Tác dụng của score.go:
- Lưu trữ điểm số do Giám khảo chấm cho từng thí sinh ở từng vòng thi.
- Hỗ trợ lưu điểm thành một map (bản đồ từ khóa) linh hoạt thay vì fix cứng cột điểm.
*/
package model

import "time"

// Score là bản ghi một lượt chấm điểm của một Giám khảo cho một Thí sinh.
type Score struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	ContestantID string `bson:"contestant_id" json:"contestantId"`
	RoundID      string `bson:"round_id" json:"roundId"`
	ExaminerID   string `bson:"examiner_id" json:"examinerId"`

	SBD            string             `bson:"sbd" json:"sbd"`
	TotalScore     float64            `bson:"total_score" json:"totalScore"`
	CriteriaScores map[string]float64 `bson:"criteria_scores" json:"criteriaScores"` // e.g., "face": 9.0, "body": 8.5
	Comment        string             `bson:"comment" json:"comment"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`

	// Audit fields (Trường kiểm toán): Dùng để truy vết xem giám khảo dùng mạng nào, máy nào để chấm (chống gian lận).
	IP        string `bson:"ip" json:"ip"`
	UserAgent string `bson:"user_agent" json:"userAgent"`
}
