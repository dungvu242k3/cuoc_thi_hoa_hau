/*
Package server chứa các tiện ích cốt lõi vận hành bên dưới ứng dụng.
Tác dụng của bcrypt.go:
  - Cung cấp công cụ mã hóa mật khẩu (hashing) cực mạnh chống lưu plain-text vào DB.
  - Tách biệt logic băm mật khẩu ra khỏi Tầng Service để nếu sau này đổi sang Argon2 hay Scrypt,
    Service không cần sửa một dòng code nào.
*/
package server

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Hash tạo ra chuỗi băm 60 ký tự từ mật khẩu gốc cộng thêm một đoạn "muối" (salt) ngẫu nhiên.
// bcrypt.DefaultCost (10) đảm bảo tốc độ tạo đủ chậm để chống brute-force nhưng đủ nhanh cho người dùng.
func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare đối chiếu mật khẩu người dùng nhập vào lúc đăng nhập với chuỗi băm trong DB.
// Do bcrypt tự nhúng salt vào trong chuỗi băm nên ta không cần lưu salt ra một cột DB riêng.
func (h *BcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
