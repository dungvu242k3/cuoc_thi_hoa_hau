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

type scoreRepo struct {
	collection *mongo.Collection
}

func NewScoreRepo(db *mongo.Database) port.ScoreRepository {
	return &scoreRepo{
		collection: db.Collection("scores"),
	}
}

func (r *scoreRepo) GetListByCandidateID(ctx context.Context, candidateID string) ([]*domain.Score, error) {
	filter := bson.M{"contestant_id": candidateID}

	// Create a cursor for querying
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []*domain.Score
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

func (r *scoreRepo) GetListByExaminerID(ctx context.Context, examinerID string) ([]*domain.Score, error) {
	filter := bson.M{"examiner_id": examinerID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []*domain.Score
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

func (r *scoreRepo) Upsert(ctx context.Context, score *domain.Score) error {
	filter := bson.M{
		"contestant_id": score.ContestantID,
		"examiner_id":   score.ExaminerID,
	}
	update := bson.M{
		"$set": score,
	}
	opts := options.Update().SetUpsert(true)

	if score.CreatedAt.IsZero() {
		score.CreatedAt = time.Now()
	}
	score.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
