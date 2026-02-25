/*
Package service đảm nhận Business Logic.
Tác dụng của schedule.go:
- Chịu trách nhiệm cung cấp Lịch trình cuộc thi cho phần Frontend.
- Kiểm soát phân trang ở tầng Nghiệp Vụ (Service), vì DAO chỉ đơn thuần là thợ đi lấy dữ liệu.
*/
package service

import (
	"context"

	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"
)

type scheduleService struct {
	repo types.ScheduleRepository
}

func NewScheduleService(repo types.ScheduleRepository) types.ScheduleService {
	return &scheduleService{repo: repo}
}

// GetPublicSchedules giới hạn query tối đa 50 item mỗi trang để chống sập server, kể cả khi Frontend truyền limit = 9999.
func (s *scheduleService) GetPublicSchedules(ctx context.Context, limit, offset int64) ([]*model.Schedule, int64, error) {
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
