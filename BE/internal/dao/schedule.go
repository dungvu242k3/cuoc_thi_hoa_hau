/*
Package dao (Data Access Object) hay còn gọi là Repository Layer.
Tác dụng của schedule.go:
- Giao tiếp trực tiếp với MongoDB (collection "schedules").
- Ẩn giấu các lệnh bson.M, Find, InsertOne của Mongo Driver khỏi tầng Service.
- Nếu sau này đổi sang PostgreSQL, chỉ cần viết lại file này, Service không bị ảnh hưởng.
*/
package dao

import (
	"context"
	"time"

	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoScheduleRepo struct {
	coll *mongo.Collection
}

func NewScheduleRepo(db *mongo.Database) types.ScheduleRepository {
	return &mongoScheduleRepo{coll: db.Collection("schedules")}
}

// GetList truy vấn danh sách lịch trình có phân trang và sắp xếp theo thời gian tăng dần.
func (r *mongoScheduleRepo) GetList(ctx context.Context, limit, offset int64) ([]*model.Schedule, int64, error) {
	filter := bson.M{}

	total, err := r.coll.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"start_time": 1}).
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*model.Schedule
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *mongoScheduleRepo) Create(ctx context.Context, s *model.Schedule) error {
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	s.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, s)
	return err
}
