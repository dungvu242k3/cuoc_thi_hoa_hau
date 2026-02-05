package domain

import (
	"testing"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		permission string
		want       bool
	}{
		{
			name:       "Admin has all permissions",
			role:       string(RoleAdmin),
			permission: PermUserWrite,
			want:       true,
		},
		{
			name:       "Examiner has score write permission",
			role:       string(RoleExaminer),
			permission: PermScoreWrite,
			want:       true,
		},
		{
			name:       "Examiner does not have user write permission",
			role:       string(RoleExaminer),
			permission: PermUserWrite,
			want:       false,
		},
		{
			name:       "Candidate has contestant read permission",
			role:       string(RoleCandidate),
			permission: PermContestantRead,
			want:       true,
		},
		{
			name:       "Candidate does not have score write permission",
			role:       string(RoleCandidate),
			permission: PermScoreWrite,
			want:       false,
		},
		{
			name:       "Invalid role has no permissions",
			role:       "invalid_role",
			permission: PermContestantRead,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.role, tt.permission); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
