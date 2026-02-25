/*
Package model định nghĩa thực thể Tài khoản người dùng (User) và Quyền hạn (Role).
Tác dụng:
- Quản lý định dạng dữ liệu lưu trong collection "users" của MongoDB.
- Phân biệt các loại quyền bằng chuỗi enum tĩnh (RoleAdmin, RoleCandidate...).
*/
package model

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleCandidate Role = "candidate"
	RoleExaminer  Role = "examiner"
	RoleAudience  Role = "audience"
)

// User là tài khoản đăng nhập để vào hệ thống.
// Một User có thể chỉ là Quản trị viên, hoặc là 1 Thí sinh (khi đó nó sẽ link với struct Contestant qua UserID).
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"password"` // Cột này chỉ được lưu chuỗi đã băm (hash), tuyệt đối KHÔNG lưu chữ thường.
	RoleID    string    `bson:"role_id" json:"roleId"`    // Changed from Role to RoleID (FK)
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// AuthClaims là thông tin được giải mã (bóc ngoặc) từ chuỗi JWT Token do user gửi lên.
// Nó chứa UserID và Role để phía Server biết "À, anh này là Giám khảo, mã ID là 123".
type AuthClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

const (
	ClaimKeyUserID = "user_id"
	ClaimKeyRole   = "role"
	ClaimKeyExp    = "exp"
)
