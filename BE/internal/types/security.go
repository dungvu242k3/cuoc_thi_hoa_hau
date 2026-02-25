/*
Package types chứa Interface cho các thuật toán bảo mật.
Tác dụng của security.go:
- PasswordEncoder: Quy định hàm mã hóa mật khẩu (ẩn giấu thư viện bcrypt).
- TokenProvider: Quy định hàm sinh chữ ký số JWT chống giả mạo danh tính.
*/
package types

import (
	"cuoc_thi_hoa_hau/internal/model"
	"time"
)

type PasswordEncoder interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type TokenProvider interface {
	Generate(user *model.User, duration time.Duration) (*model.AuthClaims, string, error)
	Validate(tokenString string) (*model.AuthClaims, error)
}
