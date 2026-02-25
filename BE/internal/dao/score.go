/*
Package dao điều hướng dữ liệu MongoDB.
Tác dụng của score.go:
- Repo chuyên dụng cho collection "scores".
- Dùng cú pháp SetUpsert (Update or Insert) để tránh trùng lặp điểm số khi Giám khảo sửa điểm nhiều lần.
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

type scoreRepo struct {
	collection *mongo.Collection
}

func NewScoreRepo(db *mongo.Database) types.ScoreRepository {
	return &scoreRepo{
		collection: db.Collection("scores"),
	}
}

func (r *scoreRepo) GetListByCandidateID(ctx context.Context, candidateID string) ([]*model.Score, error) {
	filter := bson.M{"contestant_id": candidateID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []*model.Score
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

func (r *scoreRepo) GetListByExaminerID(ctx context.Context, examinerID string) ([]*model.Score, error) {
	filter := bson.M{"examiner_id": examinerID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []*model.Score
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

// Upsert: Nếu chưa có điểm (Contestant + Examiner mapping) thì tạo mới. Nếu có rồi thì ghi đè (Update).
// Đây là Pattern "Idempotent" rất quan trọng trong thiết kế DB.
func (r *scoreRepo) Upsert(ctx context.Context, score *model.Score) error {
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
