package domain

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleCandidate Role = "candidate"
	RoleExaminer  Role = "examiner"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"password"`
	Role      Role      `bson:"role" json:"role"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

type AuthClaims struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
}
