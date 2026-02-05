package mongodb

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContestantRepo struct {
	coll *mongo.Collection
}

func NewContestantRepo(db *mongo.Database) *ContestantRepo {
	coll := db.Collection("contestants")

	// Ensure Indexes for Scaling
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		indices := []mongo.IndexModel{
			{Keys: bson.D{{Key: "is_public", Value: 1}}},
			{Keys: bson.D{{Key: "sbd", Value: 1}, {Key: "is_public", Value: 1}}},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: -1}}},
		}
		coll.Indexes().CreateMany(ctx, indices)
	}()

	return &ContestantRepo{
		coll: coll,
	}
}

func (r *ContestantRepo) Create(ctx context.Context, c *domain.Contestant) error {
	_, err := r.coll.InsertOne(ctx, c)
	return err
}

func (r *ContestantRepo) Update(ctx context.Context, c *domain.Contestant) error {
	oid, err := primitive.ObjectIDFromHex(c.ID)
	if err != nil {
		log.Printf("[Repo Error] Invalid ID hex: %v", err)
		return err
	}

	// Convert struct to bson.M to exclude _id (immutable)
	var updateDoc bson.M
	data, _ := bson.Marshal(c)
	if err := bson.Unmarshal(data, &updateDoc); err != nil {
		log.Printf("[Repo Error] Marshal failed: %v", err)
		return err
	}
	delete(updateDoc, "_id")

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": updateDoc}
	_, err = r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("[Repo Error] UpdateOne failed: %v", err)
	}
	return err
}

func (r *ContestantRepo) GetByID(ctx context.Context, id string) (*domain.Contestant, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c domain.Contestant
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&c)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ContestantRepo) GetPublicDetail(ctx context.Context, id string) (*domain.Contestant, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c domain.Contestant
	filter := bson.M{
		"_id":       oid,
		"is_public": true,
	}
	err = r.coll.FindOne(ctx, filter).Decode(&c)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ContestantRepo) GetByUserID(ctx context.Context, userId string) (*domain.Contestant, error) {
	var c domain.Contestant
	err := r.coll.FindOne(ctx, bson.M{"user_id": userId}).Decode(&c)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ContestantRepo) GetPublicList(ctx context.Context, limit, offset int64) ([]*domain.Contestant, int64, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	filter := bson.M{"is_public": true} // Only public profiles

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Contestant
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, _ := r.coll.CountDocuments(ctx, filter)
	return results, total, nil
}

func (r *ContestantRepo) CheckIdentifyCard(ctx context.Context, cardID string) (bool, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"personal_info.identity_card": cardID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ContestantRepo) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

func (r *ContestantRepo) GetList(ctx context.Context, limit, offset int64, filter map[string]interface{}) ([]*domain.Contestant, int64, error) {
	// Convert filter to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})

	cursor, err := r.coll.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Contestant
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, _ := r.coll.CountDocuments(ctx, bsonFilter)
	return results, total, nil
}

func (r *ContestantRepo) IncrementVote(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$inc": bson.M{"vote_count": 1}}
	_, err = r.coll.UpdateOne(ctx, filter, update)
	return err
}

func (r *ContestantRepo) HasVoted(ctx context.Context, userID, contestantID string) (bool, error) {
	filter := bson.M{
		"user_id":       userID,
		"contestant_id": contestantID,
	}
	count, err := r.coll.Database().Collection("vote_histories").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ContestantRepo) CheckIPLimit(ctx context.Context, ip string, contestantID string) (bool, error) {
	// Daily limit: Only count votes in the last 24 hours
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	filter := bson.M{
		"ip_address":    ip,
		"contestant_id": contestantID,
		"created_at":    bson.M{"$gte": oneDayAgo},
	}
	count, err := r.coll.Database().Collection("vote_histories").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (r *ContestantRepo) RecordVote(ctx context.Context, userID, contestantID, ip, userAgent string) error {
	doc := bson.M{
		"user_id":       userID,
		"contestant_id": contestantID,
		"ip_address":    ip,
		"user_agent":    userAgent,
		"created_at":    time.Now(),
	}
	// Create index if not exists to ensure uniqueness (best effort, ideally done once at startup)
	// For now, we rely on application logic + unique index
	_, err := r.coll.Database().Collection("vote_histories").InsertOne(ctx, doc)
	return err
}
