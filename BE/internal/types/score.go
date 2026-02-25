/*
Package types chứa Interface giao tiếp.
Tác dụng của score.go:
- ScoreRepository: Hợp đồng để tầng DAO lưu điểm thi xuống DB.
- ScoreService: Hợp đồng để tầng Service tính toán tổng điểm, chia trung bình.
*/
package types

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
)

type ScoreRepository interface {
	GetListByCandidateID(ctx context.Context, candidateID string) ([]*model.Score, error)
	GetListByExaminerID(ctx context.Context, examinerID string) ([]*model.Score, error)
	Upsert(ctx context.Context, score *model.Score) error
}

type ScoreService interface {
	GetMyScores(ctx context.Context, userID string) ([]*model.Score, error)
	SubmitScore(ctx context.Context, examinerID string, score *model.Score, ip, ua string) error
}
