/*
Package dao (Data Access Object) hay còn gọi là Repository Layer.
Tác dụng:
  - Nơi DUY NHẤT trong toàn bộ dự án trực tiếp giao tiếp với Database (ở đây là MongoDB).
  - Chứa các thư viện như go.mongodb.org. Bất kỳ tầng nào khác (Handler, Service) có import thư viện
    MongoDB đều là vi phạm kiến trúc.
  - Ẩn giấu sự phức tạp của DB đi, chỉ cung cấp các hàm Insert, Find sạch sẽ cho Tầng Service gọi.
*/
package dao

import (
	"context"
	"time"

	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) types.UserRepository {
	repo := &UserRepo{db: db}
	repo.ensureIndex()
	return repo
}

func (r *UserRepo) ensureIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
	}
}

// Create thực thi lệnh InsertOne thẳng xuống MongoDB Collection "users".
// Nó sẽ tự động sinh CreatedAt, UpdatedAt và ID nếu DB thành công.
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	res, err := r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}
	return nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.db.Collection("users").FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
