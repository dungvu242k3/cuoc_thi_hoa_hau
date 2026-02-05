package resolver

import (
	"cuoc_thi_hoa_hau/internal/adapter/graph/model"
	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/pkg/common"
	"time"
)

// mapContestantToGraphFull returns full data (for Admin/Owner only)
func mapContestantToGraphFull(d *domain.Contestant) *model.Contestant {
	if d == nil {
		return nil
	}
	return &model.Contestant{
		ID:       d.ID,
		UserID:   d.UserID,
		Sbd:      &d.SBD,
		Status:   string(d.Status),
		IsPublic: d.IsPublic,
		PersonalInfo: &model.PersonalInfo{
			FullName:     d.PersonalInfo.FullName,
			Dob:          d.PersonalInfo.DateOfBirth,
			Nationality:  &d.PersonalInfo.Nationality,
			Gender:       &d.PersonalInfo.Gender,
			IdentityCard: &d.PersonalInfo.IdentityCard,
			Phone:        &d.PersonalInfo.Phone,
			Email:        &d.PersonalInfo.Email,
			Address:      &d.PersonalInfo.Address,
			Job:          &d.PersonalInfo.Job,
		},
		PhysicalInfo: &model.PhysicalInfo{
			Height:       d.PhysicalInfo.Height,
			Weight:       d.PhysicalInfo.Weight,
			Measurements: &d.PhysicalInfo.Measurements,
		},
		SkillEducation: &model.SkillEducation{
			EducationLevel: &d.SkillEdu.EducationLevel,
			Languages:      d.SkillEdu.Languages,
			Skills:         d.SkillEdu.Skills,
		},
		Portfolio: &model.Portfolio{
			AvatarURL:    &d.Portfolio.AvatarURL,
			GalleryUrls:  d.Portfolio.GalleryURLs,
			VideoURL:     &d.Portfolio.VideoURL,
			Introduction: &d.Portfolio.Introduction,
			SocialLinks:  d.Portfolio.SocialLinks,
		},
		VoteCount: int(d.VoteCount),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// mapContestantToGraphOwner returns full data for owner (contestant viewing their own profile)
func mapContestantToGraphOwner(d *domain.Contestant) *model.Contestant {
	return mapContestantToGraphFull(d)
}

func mapPublicContestantToGraph(d *domain.Contestant) *model.Contestant {
	if d == nil {
		return nil
	}

	// Masking Sensitive Data
	masked := "******"

	res := mapContestantToGraphFull(d)
	res.PersonalInfo.IdentityCard = &masked
	res.PersonalInfo.Phone = &masked
	res.PersonalInfo.Email = &masked

	return res
}

func ToDomainContestant(input model.CreateContestantInput, dob time.Time) *domain.Contestant {
	return &domain.Contestant{
		PersonalInfo: domain.PersonalInfo{
			FullName:     input.FullName,
			DateOfBirth:  dob,
			Nationality:  input.Nationality,
			Gender:       input.Gender,
			IdentityCard: getString(input.IdentityCard),
			Phone:        input.Phone,
			Email:        input.Email,
			Address:      input.Address,
			Job:          input.Job,
		},
		PhysicalInfo: domain.PhysicalInfo{
			Height:       input.Height,
			Weight:       input.Weight,
			Measurements: input.Measurements,
		},
		SkillEdu: domain.SkillEducation{
			EducationLevel: getString(input.EducationLevel),
			Languages:      input.Languages,
			Skills:         input.Skills,
		},
		Portfolio: domain.Portfolio{
			AvatarURL:    getString(input.AvatarURL),
			GalleryURLs:  input.GalleryUrls,
			Introduction: getString(input.Introduction),
			SocialLinks:  input.SocialLinks,
		},
	}
}

func ToDomainUpdateContestant(input model.UpdateContestantInput) (*domain.Contestant, error) {
	d := &domain.Contestant{}
	if input.FullName != nil {
		d.PersonalInfo.FullName = *input.FullName
	}
	if input.Dob != nil {
		dob, err := time.Parse(common.DateLayoutISO, *input.Dob)
		if err != nil {
			return nil, err
		}
		d.PersonalInfo.DateOfBirth = dob
	}
	if input.Nationality != nil {
		d.PersonalInfo.Nationality = *input.Nationality
	}
	if input.Gender != nil {
		d.PersonalInfo.Gender = *input.Gender
	}
	if input.Phone != nil {
		d.PersonalInfo.Phone = *input.Phone
	}
	if input.Email != nil {
		d.PersonalInfo.Email = *input.Email
	}
	if input.Address != nil {
		d.PersonalInfo.Address = *input.Address
	}
	if input.Job != nil {
		d.PersonalInfo.Job = *input.Job
	}

	if input.Height != nil {
		d.PhysicalInfo.Height = *input.Height
	}
	if input.Weight != nil {
		d.PhysicalInfo.Weight = *input.Weight
	}
	if input.Measurements != nil {
		d.PhysicalInfo.Measurements = *input.Measurements
	}

	if input.EducationLevel != nil {
		d.SkillEdu.EducationLevel = *input.EducationLevel
	}
	if input.Languages != nil {
		d.SkillEdu.Languages = input.Languages
	}
	if input.Skills != nil {
		d.SkillEdu.Skills = input.Skills
	}

	if input.AvatarURL != nil {
		d.Portfolio.AvatarURL = *input.AvatarURL
	}
	if input.GalleryUrls != nil {
		d.Portfolio.GalleryURLs = input.GalleryUrls
	}
	if input.Introduction != nil {
		d.Portfolio.Introduction = *input.Introduction
	}
	if input.SocialLinks != nil {
		d.Portfolio.SocialLinks = input.SocialLinks
	}
	return d, nil
}

func ToDomainScore(input model.ScoreInput) *domain.Score {
	s := &domain.Score{
		ContestantID:   input.ContestantID,
		SBD:            input.Sbd,
		CriteriaScores: make(map[string]float64),
		Comment:        getString(input.Comment),
	}

	for k, v := range input.CriteriaScores {
		val, err := safelyToFloat64(v)
		if err == nil {
			s.CriteriaScores[k] = val
		}
	}
	return s
}

func ToGraphScore(s *domain.Score) *model.Score {
	if s == nil {
		return nil
	}

	criteria := make(map[string]any)
	for k, v := range s.CriteriaScores {
		criteria[k] = v
	}

	var commentPtr *string
	if s.Comment != "" {
		val := s.Comment
		commentPtr = &val
	}

	return &model.Score{
		ID:             s.ID,
		ContestantID:   s.ContestantID,
		RoundID:        s.RoundID,
		Sbd:            s.SBD,
		TotalScore:     s.TotalScore,
		CriteriaScores: criteria,
		Comment:        commentPtr,
		CreatedAt:      s.CreatedAt,
	}
}

// getString handles nil string pointers safely
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
