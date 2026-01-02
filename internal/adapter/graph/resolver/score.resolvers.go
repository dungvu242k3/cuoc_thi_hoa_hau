package resolver

import (
	"context"
	"errors"

	"cuoc_thi_hoa_hau/internal/adapter/graph/middleware"
	"cuoc_thi_hoa_hau/internal/adapter/graph/model"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

func (r *queryResolver) MyScore(ctx context.Context) (*model.ScoreDetail, error) {
	// Auth Check
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok {
		return nil, errors.New("unauthorized: vui lòng đăng nhập")
	}

	//Call Service
	score, err := r.ScoringService.GetMyScore(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	if score == nil {
		return nil, nil
	}

	// Map Domain -> Model
	var details []*model.ScoreCriterion
	for k, v := range score.Details {
		details = append(details, &model.ScoreCriterion{
			Key:   k,
			Value: v,
		})
	}

	return &model.ScoreDetail{
		TotalScore: score.TotalScore,
		Details:    details,
	}, nil
}
