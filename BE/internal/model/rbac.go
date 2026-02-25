package model

import "time"

type Permission struct {
	Code        string    `bson:"_id" json:"code"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
}

type RoleDef struct {
	ID          string    `bson:"_id" json:"id"` // e.g., "admin", "examiner"
	Name        string    `bson:"name" json:"name"`
	Permissions []string  `bson:"permissions" json:"permissions"` // List of Permission Codes
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

// Predefined Permissions Constants
const (
	PermScoreWrite      = "score:write"
	PermScoreRead       = "score:read"
	PermContestantRead  = "contestant:read"
	PermContestantWrite = "contestant:write"
	PermSystemConfig    = "system:config"
	PermUserRead        = "user:read"
	PermUserWrite       = "user:write"
)

// RolePermissions maps roles to their respective permission codes.
var RolePermissions = map[string][]string{
	string(RoleAdmin): {
		PermScoreWrite, PermScoreRead,
		PermContestantRead, PermContestantWrite,
		PermSystemConfig, PermUserRead, PermUserWrite,
	},
	string(RoleExaminer): {
		PermScoreWrite, PermScoreRead,
		PermContestantRead, PermContestantWrite,
	},
	string(RoleCandidate): {
		PermContestantRead,
	},
	string(RoleAudience): {
		PermContestantRead,
	},
}

func HasPermission(role string, requiredPerm string) bool {
	if role == string(RoleAdmin) {
		return true
	}
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == requiredPerm {
			return true
		}
	}
	return false
}
