/*
Package types định nghĩa hợp đồng quản lý Lịch Trình.
Tác dụng:
- Tách biệt logic lấy danh sách sự kiện (Service) và câu lệnh query DB (Repository).
*/
package types

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
)

type ScheduleRepository interface {
	GetList(ctx context.Context, limit int64, offset int64) ([]*model.Schedule, int64, error)
	Create(ctx context.Context, s *model.Schedule) error
}

type ScheduleService interface {
	GetPublicSchedules(ctx context.Context, limit int64, offset int64) ([]*model.Schedule, int64, error)
}
