package resolver

import (
	"context"
	"time"

	"cuoc_thi_hoa_hau/internal/adapter/graph/model"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

func (r *queryResolver) PublicSchedules(ctx context.Context, limit *int, offset *int) (*model.ScheduleList, error) {
	// Standardize input
	l := int64(10)
	if limit != nil {
		l = int64(*limit)
	}
	o := int64(0)
	if offset != nil {
		o = int64(*offset)
	}

	// Call Service
	schedules, total, err := r.ScheduleService.GetPublicSchedules(ctx, l, o)
	if err != nil {
		return nil, err
	}

	// Map Domain -> GraphQL Model
	items := make([]*model.Schedule, 0)
	for _, s := range schedules {
		items = append(items, mapToScheduleGraphModel(s))
	}

	// Return Pagination Wrapper
	return &model.ScheduleList{
		Items: items,
		Total: int(total),
	}, nil
}

func mapToScheduleGraphModel(s *domain.Schedule) *model.Schedule {
	if s == nil {
		return nil
	}
	return &model.Schedule{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		Type:        model.ScheduleType(s.Type),
		Location:    s.Location,
		StartTime:   s.StartTime.Format(time.RFC3339),
		EndTime:     s.EndTime.Format(time.RFC3339),
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
	}
}
