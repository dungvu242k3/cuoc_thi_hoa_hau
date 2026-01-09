package port

import (
	"context"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

// lưu lấy data ở db
type ScheduleRepository interface {
	// lấy danh sách
	GetList(ctx context.Context, limit int64, offset int64) ([]*domain.Schedule, int64, error)
	//tạo lịch trình mới
	Create(ctx context.Context, s *domain.Schedule) error
}

type ScheduleService interface {
	// cung cấp danh sách cho thí sinh
	GetPublicSchedules(ctx context.Context, limit int64, offset int64) ([]*domain.Schedule, int64, error)
}
