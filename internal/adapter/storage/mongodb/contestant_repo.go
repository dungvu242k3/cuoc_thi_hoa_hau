package mongodb

import (
	"context"
	"time"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	coll *mongo.Collection
}

func NewContestantRepo(db *mongo.Database) port.ContestantRepository {
	return &mongoRepo{coll: db.Collection("contestants")}
}

// Projection để bảo mật thông tin
// Projection để bảo mật thông tin (Defense in Depth layer 1)
var publicProjection = bson.M{
	"_id": 1, "status": 1, "is_public": 1, "created_at": 1,
	"personal_info.full_name": 1, "personal_info.dob": 1, "personal_info.nationality": 1,
	// KHÔNG CÓ: Phone, Email, Address, IdentifyCard
	"physical_info":   1,
	"portfolio":       1,
	"skill_education": 1,
}

func (r *mongoRepo) GetPublicList(ctx context.Context, limit, offset int64) ([]*domain.Contestant, int64, error) {
	filter := bson.M{"status": domain.StatusApproved, "is_public": true}

	// 1. Count Total
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 2. Query Data
	opts := options.Find().
		SetProjection(publicProjection).
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var results []*domain.Contestant
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *mongoRepo) GetPublicDetail(ctx context.Context, id string) (*domain.Contestant, error) {
	filter := bson.M{"_id": id, "status": domain.StatusApproved, "is_public": true}

	var result domain.Contestant
	err := r.coll.FindOne(ctx, filter, options.FindOne().SetProjection(publicProjection)).Decode(&result)
	return &result, err
}

func (r *mongoRepo) Create(ctx context.Context, c *domain.Contestant) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, c)
	return err
}

func (r *mongoRepo) Update(ctx context.Context, c *domain.Contestant) error {
	c.UpdatedAt = time.Now()

	// Filter: Support both ObjectID and String
	filter := bson.M{"_id": c.ID}
	if objID, err := primitive.ObjectIDFromHex(c.ID); err == nil {
		filter = bson.M{"_id": objID}
	}

	// Update payload: Strip _id to avoid "immutable field" error
	updateData := bson.M{}
	doc, _ := bson.Marshal(c)
	bson.Unmarshal(doc, &updateData)
	delete(updateData, "_id")

	_, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": updateData})
	return err
}

func (r *mongoRepo) GetByUserID(ctx context.Context, uid string) (*domain.Contestant, error) {
	var res domain.Contestant
	err := r.coll.FindOne(ctx, bson.M{"user_id": uid}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *mongoRepo) GetByID(ctx context.Context, id string) (*domain.Contestant, error) {
	var result domain.Contestant
	filter := bson.M{"_id": id}
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objID}
	}
	err := r.coll.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found
	}
	return &result, err
}

func (r *mongoRepo) CheckIdentifyCard(ctx context.Context, card string) (bool, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"identify_card": card})
	return count > 0, err
}
