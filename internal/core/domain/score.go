package domain

import "time"

type Score struct {
	ID          string             `bson:"_id,omitempty" json:"id"`
	CandidateID string             `bson:"candidate_id" json:"candidateId"`
	SBD         string             `bson:"sbd" json:"sbd"`
	TotalScore  float64            `bson:"total_score" json:"totalScore"`
	Details     map[string]float64 `bson:"details" json:"details"`
	ExaminerID  string             `bson:"examiner_id,omitempty" json:"-"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}
