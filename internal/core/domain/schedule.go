package domain

import "time"

type ScheduleType string

const (
	ScheduleExam     ScheduleType = "EXAM"     // Lịch thi
	ScheduleTraining ScheduleType = "TRAINING" // Lịch tập luyện
	ScheduleEvent    ScheduleType = "EVENT"    // Sự kiện truyền thông
)

type Schedule struct {
	ID          string       `bson:"_id,omitempty" json:"id"`
	Title       string       `bson:"title" json:"title"`
	Description string       `bson:"description" json:"description"`
	Type        ScheduleType `bson:"type" json:"type"`
	Location    string       `bson:"location" json:"location"`

	StartTime time.Time `bson:"start_time" json:"startTime"`
	EndTime   time.Time `bson:"end_time" json:"endTime"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
