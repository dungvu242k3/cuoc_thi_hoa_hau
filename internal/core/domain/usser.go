package domain

type AuthClaims struct {
	UserID string
	Role   string
}

const (
	RoleCandidate = "CANDIDATE"
	RoleAdmin     = "ADMIN"
)
