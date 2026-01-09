package domain

import "time"

type ScheduleType string

const (
	ScheduleTypeExam     ScheduleType = "exam"
	ScheduleTypeTraining ScheduleType = "training"
	ScheduleTypeMedia    ScheduleType = "media"
	ScheduleTypeOther    ScheduleType = "other"
)

type Schedule struct {
	ID          string       `bson:"_id,omitempty" json:"id"`
	Title       string       `bson:"title" json:"title"`
	Description string       `bson:"description" json:"description"`
	StartTime   time.Time    `bson:"start_time" json:"startTime"`
	EndTime     time.Time    `bson:"end_time" json:"endTime"`
	Location    string       `bson:"location" json:"location"`
	Type        ScheduleType `bson:"type" json:"type"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
