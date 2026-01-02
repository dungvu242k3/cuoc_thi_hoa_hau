package resolver

import (
	"context"
	"errors"
	"time"

	"cuoc_thi_hoa_hau/internal/adapter/graph/middleware"
	"cuoc_thi_hoa_hau/internal/adapter/graph/model"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

// SendFeedback allows a candidate to submit a feedback.
func (r *mutationResolver) SendFeedback(ctx context.Context, input model.SendFeedbackInput) (bool, error) {
	// 1. Auth Check: Only Candidates can complain/propose
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok || user.Role != domain.RoleCandidate {
		return false, errors.New("bạn không có quyền thực hiện hành động này")
	}

	// 2. Map Input -> Domain
	feedback := &domain.Feedback{
		Title:   input.Title,
		Content: input.Content,
		Type:    domain.FeedbackType(input.Type),
	}

	// 3. Call Service
	if err := r.FeedbackService.SendFeedback(ctx, user.UserID, feedback); err != nil {
		return false, err
	}

	return true, nil
}

// MyFeedbacks retrieves feedbacks submitted by the requesting user.
func (r *queryResolver) MyFeedbacks(ctx context.Context, limit *int, offset *int) (*model.FeedbackList, error) {
	// 1. Auth Check
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	// 2. Standardize Pagination
	l := int64(10)
	if limit != nil {
		l = int64(*limit)
	}
	o := int64(0)
	if offset != nil {
		o = int64(*offset)
	}

	// 3. Call Service
	list, total, err := r.FeedbackService.GetMyFeedbacks(ctx, user.UserID, l, o)
	if err != nil {
		return nil, err
	}

	// 4. Map Domain -> Model
	items := make([]*model.Feedback, 0)
	for _, f := range list {
		items = append(items, &model.Feedback{
			ID:        f.ID,
			UserID:    f.UserID,
			Title:     f.Title,
			Content:   f.Content,
			Type:      model.FeedbackType(f.Type),
			Status:    model.FeedbackStatus(f.Status),
			CreatedAt: f.CreatedAt.Format(time.RFC3339),
		})
	}

	return &model.FeedbackList{
		Items: items,
		Total: int(total),
	}, nil
}
