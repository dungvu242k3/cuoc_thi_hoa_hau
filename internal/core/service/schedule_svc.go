package service

import (
	"context"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
)

type scheduleService struct {
	repo port.ScheduleRepository
}

func NewScheduleService(repo port.ScheduleRepository) port.ScheduleService {
	return &scheduleService{repo: repo}
}

func (s *scheduleService) GetPublicSchedules(ctx context.Context, limit, offset int64) ([]*domain.Schedule, int64, error) {
	// 1. Validation & Optimization: Limit max items to prevent DoS
	const maxLimit = 50
	const defaultLimit = 10

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetList(ctx, limit, offset)
}
