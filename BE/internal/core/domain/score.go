package domain

import "time"

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

	// Audit fields
	IP        string `bson:"ip" json:"ip"`
	UserAgent string `bson:"user_agent" json:"userAgent"`
}
