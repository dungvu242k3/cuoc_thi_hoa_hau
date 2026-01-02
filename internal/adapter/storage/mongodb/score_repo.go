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

func (r *scoreRepo) GetByCandidateID(ctx context.Context, candidateID string) (*domain.Score, error) {
	filter := bson.M{"candidate_id": candidateID}
	var score domain.Score
	err := r.collection.FindOne(ctx, filter).Decode(&score)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no score found (not an error)
		}
		return nil, err
	}
	return &score, nil
}

func (r *scoreRepo) Upsert(ctx context.Context, score *domain.Score) error {
	filter := bson.M{"candidate_id": score.CandidateID}
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
