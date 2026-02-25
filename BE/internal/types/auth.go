/*
Package types (Interface Layer) định nghĩa các hợp đồng giao tiếp.
Tác dụng của auth.go:
- Quy định Tầng Service của Auth bắt buộc phải triển khai hàm Register và Login.
- Handler sẽ gọi interface này thay vì gọi trực tiếp struct AuthService.
*/
package types

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, username, password, role string) (*model.AuthClaims, string, error)
	Login(ctx context.Context, username, password string) (*model.AuthClaims, string, error)
}
