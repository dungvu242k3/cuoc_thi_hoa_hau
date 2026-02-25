/*
Package resolver xử lý việc điều hướng requests từ GraphQL đến Service.
Tác dụng của mapper.go:
- Đóng vai trò Data Transfer Object (DTO) Mapper.
- Nó dịch các struct DB (trong gói internal/model) sang struct GraphQL (trong gói graph/model) trả về cho Client.
- Đồng thời dịch chiều ngược lại: Input từ Frontend thành Struct cho DB.
- Nơi đây lý tưởng để che giấu các trường nhạy cảm (như che dấu CCCD khi là Public Contestant).
*/
package resolver

import (
	gqlmodel "cuoc_thi_hoa_hau/graph/model"
	"cuoc_thi_hoa_hau/internal/model"
	"time"
)

func mapContestantToGraphFull(d *model.Contestant) *gqlmodel.Contestant {
	if d == nil {
		return nil
	}
	return &gqlmodel.Contestant{
		ID:       d.ID,
		UserID:   d.UserID,
		Sbd:      &d.SBD,
		Status:   string(d.Status),
		IsPublic: d.IsPublic,
		PersonalInfo: &gqlmodel.PersonalInfo{
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
		PhysicalInfo: &gqlmodel.PhysicalInfo{
			Height:       d.PhysicalInfo.Height,
			Weight:       d.PhysicalInfo.Weight,
			Measurements: &d.PhysicalInfo.Measurements,
		},
		SkillEducation: &gqlmodel.SkillEducation{
			EducationLevel: &d.SkillEdu.EducationLevel,
			Languages:      d.SkillEdu.Languages,
			Skills:         d.SkillEdu.Skills,
		},
		Portfolio: &gqlmodel.Portfolio{
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

func mapContestantToGraphOwner(d *model.Contestant) *gqlmodel.Contestant {
	return mapContestantToGraphFull(d)
}

// mapPublicContestantToGraph dịch dữ liệu sang GraphQL nhưng HÌNH PHẠT (MASK) làm ẩn đi số điện thoại, email, CCCD
// nếu người gọi query là khán giả bình thường (không phải chủ nhân hồ sơ hay admin).
func mapPublicContestantToGraph(d *model.Contestant) *gqlmodel.Contestant {
	if d == nil {
		return nil
	}

	masked := "******"

	res := mapContestantToGraphFull(d)
	res.PersonalInfo.IdentityCard = &masked
	res.PersonalInfo.Phone = &masked
	res.PersonalInfo.Email = &masked

	return res
}

func ToDomainContestant(input gqlmodel.CreateContestantInput, dob time.Time) *model.Contestant {
	return &model.Contestant{
		PersonalInfo: model.PersonalInfo{
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
		PhysicalInfo: model.PhysicalInfo{
			Height:       input.Height,
			Weight:       input.Weight,
			Measurements: input.Measurements,
		},
		SkillEdu: model.SkillEducation{
			EducationLevel: getString(input.EducationLevel),
			Languages:      input.Languages,
			Skills:         input.Skills,
		},
		Portfolio: model.Portfolio{
			AvatarURL:    getString(input.AvatarURL),
			GalleryURLs:  input.GalleryUrls,
			Introduction: getString(input.Introduction),
			SocialLinks:  input.SocialLinks,
		},
	}
}

func ToDomainUpdateContestant(input gqlmodel.UpdateContestantInput) (*model.Contestant, error) {
	d := &model.Contestant{}
	if input.FullName != nil {
		d.PersonalInfo.FullName = *input.FullName
	}
	if input.Dob != nil {
		dob, err := time.Parse(model.DateLayoutISO, *input.Dob)
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

func ToDomainScore(input gqlmodel.ScoreInput) *model.Score {
	s := &model.Score{
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

func ToGraphScore(s *model.Score) *gqlmodel.Score {
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

	return &gqlmodel.Score{
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

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
