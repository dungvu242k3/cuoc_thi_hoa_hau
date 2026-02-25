/*
Package types (hoặc Ports trong Cấu trúc Hình lục giác/Hexagonal).
Tác dụng: Định nghĩa sẵn "Bản Hợp Đồng" (Interface) giữa các tầng.
  - Nhờ khai báo Interface ở thư mục trung lập này, Tầng Service có thể gọi Tầng DAO
    mà không cần phải import trực tiếp package dao (tránh lỗi vòng tròn import cycle).
  - Rất dễ viết Unit Test vì có thể tạo MockRepo giả lập interface này.
*/
package types

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}
