/*
Package model chia sẻ cấu trúc Lịch Trình (Schedules).
Tác dụng của schedule.go:
- Lưu trữ các sự kiện, lịch tập luyện, thi đấu của cuộc thi (lưu ở collection schedules).
*/
package model

import "time"

type ScheduleType string

const (
	ScheduleTypeExam     ScheduleType = "exam"
	ScheduleTypeTraining ScheduleType = "training"
	ScheduleTypeMedia    ScheduleType = "media"
	ScheduleTypeOther    ScheduleType = "other"
)

// Schedule là bản ghi lưu trữ một sự kiện trong khuôn khổ cuộc thi.
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
