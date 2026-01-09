package mongodb

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContestantRepo struct {
	coll *mongo.Collection
}

func NewContestantRepo(db *mongo.Database) *ContestantRepo {
	return &ContestantRepo{
		coll: db.Collection("contestants"),
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
