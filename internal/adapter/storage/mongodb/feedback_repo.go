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

type mongoFeedbackRepo struct {
	coll *mongo.Collection
}

func NewFeedbackRepo(db *mongo.Database) port.FeedbackRepository {
	return &mongoFeedbackRepo{coll: db.Collection("feedbacks")}
}

func (r *mongoFeedbackRepo) Create(ctx context.Context, f *domain.Feedback) error {
	now := time.Now()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	f.UpdatedAt = now
	// Default status if not set
	if f.Status == "" {
		f.Status = domain.FeedbackPending
	}

	_, err := r.coll.InsertOne(ctx, f)
	return err
}

func (r *mongoFeedbackRepo) GetListByUser(ctx context.Context, userID string, limit, offset int64) ([]*domain.Feedback, int64, error) {
	filter := bson.M{"user_id": userID}

	// CountDocuments is used here because we are filtering by user_id
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}). // Mới nhất lên đầu
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Feedback
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}
	return results, total, nil
}
