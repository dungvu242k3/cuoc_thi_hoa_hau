package domain

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleCandidate Role = "candidate"
	RoleExaminer  Role = "examiner"
	RoleAudience  Role = "audience"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"password"`
	RoleID    string    `bson:"role_id" json:"roleId"` // Changed from Role to RoleID (FK)
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

type AuthClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

const (
	ClaimKeyUserID = "user_id"
	ClaimKeyRole   = "role"
	ClaimKeyExp    = "exp"
)
