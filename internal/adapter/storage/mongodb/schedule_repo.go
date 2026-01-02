package mongodb

import (
	"context"
	"time"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoScheduleRepo struct {
	coll *mongo.Collection
}

func NewScheduleRepo(db *mongo.Database) port.ScheduleRepository {
	return &mongoScheduleRepo{coll: db.Collection("schedules")}
}

func (r *mongoScheduleRepo) GetList(ctx context.Context, limit, offset int64) ([]*domain.Schedule, int64, error) {
	// Query: Lấy tất cả
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

	var results []*domain.Schedule
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *mongoScheduleRepo) Create(ctx context.Context, s *domain.Schedule) error {
	// Logic: Tự động set thời gian tạo/cập nhật
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	s.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, s)
	return err
}
