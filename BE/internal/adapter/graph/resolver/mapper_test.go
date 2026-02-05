package resolver

import (
	"cuoc_thi_hoa_hau/internal/core/domain"
	"testing"
)

func TestMapPublicContestantToGraph(t *testing.T) {
	c := &domain.Contestant{
		ID:     "1",
		UserID: "user1",
		SBD:    "001",
		PersonalInfo: domain.PersonalInfo{
			FullName:     "Test User",
			IdentityCard: "123456789",
			Phone:        "0901234567",
			Email:        "test@example.com",
		},
	}

	got := mapPublicContestantToGraph(c)

	if *got.PersonalInfo.IdentityCard != "******" {
		t.Errorf("expected IdentityCard to be masked, got %s", *got.PersonalInfo.IdentityCard)
	}
	if *got.PersonalInfo.Phone != "******" {
		t.Errorf("expected Phone to be masked, got %s", *got.PersonalInfo.Phone)
	}
	if *got.PersonalInfo.Email != "******" {
		t.Errorf("expected Email to be masked, got %s", *got.PersonalInfo.Email)
	}
	if got.PersonalInfo.FullName != "Test User" {
		t.Errorf("expected FullName to be 'Test User', got %s", got.PersonalInfo.FullName)
	}
}
