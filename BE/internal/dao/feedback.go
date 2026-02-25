/*
Package dao chứa các implementation cho giao tiếp Database.
Tác dụng của feedback.go:
- Thao tác độc quyền trên collection "feedbacks".
- Đảm bảo các Audit fields (CreatedAt, UpdatedAt) được set tự động trước khi lưu.
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

type mongoFeedbackRepo struct {
	coll *mongo.Collection
}

func NewFeedbackRepo(db *mongo.Database) types.FeedbackRepository {
	return &mongoFeedbackRepo{coll: db.Collection("feedbacks")}
}

func (r *mongoFeedbackRepo) Create(ctx context.Context, f *model.Feedback) error {
	now := time.Now()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	f.UpdatedAt = now
	if f.Status == "" {
		f.Status = model.FeedbackStatusPending
	}

	_, err := r.coll.InsertOne(ctx, f)
	return err
}

// GetListByUser lấy danh sách phản hồi của một User cụ thể, sắp xếp mới nhất lên đầu.
func (r *mongoFeedbackRepo) GetListByUser(ctx context.Context, userID string, limit, offset int64) ([]*model.Feedback, int64, error) {
	filter := bson.M{"contestant_id": userID}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*model.Feedback
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}
	return results, total, nil
}
