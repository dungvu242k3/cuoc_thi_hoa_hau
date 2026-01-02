package resolver

import (
	"context"
	"errors"
	"time"

	"cuoc_thi_hoa_hau/internal/adapter/graph/middleware"
	"cuoc_thi_hoa_hau/internal/adapter/graph/model"
	"cuoc_thi_hoa_hau/internal/core/domain"
)

// MUTATION (Ghi)

func (r *mutationResolver) CreateContestantProfile(ctx context.Context, input model.CreateContestantInput) (*model.Contestant, error) {
	// 1Auth Check
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok || user.Role != domain.RoleCandidate {
		return nil, errors.New("bạn không có quyền thực hiện hành động này")
	}

	// Map Input -> Domain Entity
	dob, _ := time.Parse("2006-01-02", input.Dob)
	entity := &domain.Contestant{
		IdentifyCard: input.IdentifyCard,
		PersonalInfo: domain.PersonalInfo{
			FullName: input.FullName, DateOfBirth: dob, Nationality: input.Nationality,
			Phone: input.Phone, Email: input.Email, Address: input.Address,
		},
		PhysicalInfo: domain.PhysicalInfo{
			Height: input.Height, Weight: input.Weight, Measurements: input.Measurements,
		},
	}

	// Call Service
	created, err := r.ContestantService.CreateProfile(ctx, user.UserID, entity)
	if err != nil {
		return nil, err
	}
	return mapToGraphModel(created), nil
}

func (r *mutationResolver) UpdateContestantProfile(ctx context.Context, input model.UpdateContestantInput) (*model.Contestant, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok || user.Role != domain.RoleCandidate {
		return nil, errors.New("unauthorized")
	}

	updateData := &domain.Contestant{}
	if input.FullName != nil {
		updateData.PersonalInfo.FullName = *input.FullName
	}
	// Map các trường còn lại tương tự

	updated, err := r.ContestantService.UpdateProfile(ctx, user.UserID, updateData)
	if err != nil {
		return nil, err
	}
	return mapToGraphModel(updated), nil
}

func (r *mutationResolver) SubmitProfile(ctx context.Context) (bool, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok || user.Role != domain.RoleCandidate {
		return false, errors.New("unauthorized")
	}
	return true, r.ContestantService.SubmitProfile(ctx, user.UserID)
}

// QUERY đọc

func (r *queryResolver) MyProfile(ctx context.Context) (*model.Contestant, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*domain.AuthClaims)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	profile, err := r.ContestantService.GetMyProfile(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	return mapToGraphModel(profile), nil
}

func (r *queryResolver) PublicContestants(ctx context.Context, page *int, limit *int) ([]*model.PublicContestant, error) {
	p, l := 1, 20
	if page != nil && *page > 0 {
		p = *page
	}
	if limit != nil && *limit > 0 && *limit <= 100 {
		l = *limit
	}
	offset := int64((p - 1) * l)

	list, total, err := r.ContestantService.GetPublicList(ctx, int64(l), offset)
	if err != nil {
		return nil, err
	}
	// Note: Total count could be returned if schema supported Pagination object
	_ = total

	var res []*model.PublicContestant
	for _, c := range list {
		// 1. Sanitize at Domain Level first (Defense in Depth)
		safe := c.ToPublicView()
		// 2. Map to GraphQL Model
		res = append(res, mapToPublicGraphModel(safe))
	}
	return res, nil
}

func (r *queryResolver) PublicContestantDetail(ctx context.Context, id string) (*model.PublicContestant, error) {
	// Note: ContestantService cần thêm hàm GetByID hoặc dùng lại Repo
	// Ở đây giả sử dùng Repo hoặc Service có hàm tương tự
	// Để đơn giản ta gọi Service.GetPublicDetail (cần thêm vào service) hoặc dùng tạm GetByID rồi sanitize
	profile, err := r.ContestantService.GetPublicDetail(ctx, id) // Giả sử đã thêm hàm này
	if err != nil {
		return nil, err
	}
	return mapToPublicGraphModel(profile.ToPublicView()), nil
}

func mapToPublicGraphModel(d *domain.Contestant) *model.PublicContestant {
	if d == nil {
		return nil
	}
	return &model.PublicContestant{
		ID:        d.ID,
		Status:    string(d.Status),
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
		PersonalInfo: &model.PublicPersonalInfo{
			FullName:    d.PersonalInfo.FullName,
			Dob:         d.PersonalInfo.DateOfBirth.Format("2006-01-02"),
			Nationality: &d.PersonalInfo.Nationality,
		},
		PhysicalInfo: &model.PhysicalInfo{
			Height:       d.PhysicalInfo.Height,
			Weight:       d.PhysicalInfo.Weight,
			Measurements: d.PhysicalInfo.Measurements,
		},
		SkillEducation: &model.SkillEducation{
			EducationLevel: &d.SkillEdu.EducationLevel,
			Languages:      convertStringSlice(d.SkillEdu.Languages),
			Skills:         convertStringSlice(d.SkillEdu.Skills),
		},
		Portfolio: &model.Portfolio{
			AvatarURL:    &d.Portfolio.AvatarURL,
			GalleryURLs:  convertStringSlice(d.Portfolio.GalleryURLs),
			Introduction: &d.Portfolio.Introduction,
			SocialLinks:  convertStringSlice(d.Portfolio.SocialLinks),
		},
	}
}

func convertStringSlice(in []string) []string {

	return in
}

func mapToGraphModel(d *domain.Contestant) *model.Contestant {
	if d == nil {
		return nil
	}
	return &model.Contestant{
		ID: d.ID, UserID: d.UserID, Status: string(d.Status), IsPublic: d.IsPublic,
		PersonalInfo: &model.PersonalInfo{
			FullName:    d.PersonalInfo.FullName,
			Dob:         d.PersonalInfo.DateOfBirth.Format("2006-01-02"),
			Phone:       &d.PersonalInfo.Phone,
			Email:       &d.PersonalInfo.Email,
			Address:     &d.PersonalInfo.Address,
			Nationality: &d.PersonalInfo.Nationality,
		},
		PhysicalInfo: &model.PhysicalInfo{
			Height:       d.PhysicalInfo.Height,
			Weight:       d.PhysicalInfo.Weight,
			Measurements: d.PhysicalInfo.Measurements,
		},
		SkillEducation: &model.SkillEducation{
			EducationLevel: &d.SkillEdu.EducationLevel,
			Languages:      convertStringSlice(d.SkillEdu.Languages),
			Skills:         convertStringSlice(d.SkillEdu.Skills),
		},
		Portfolio: &model.Portfolio{
			AvatarURL:    &d.Portfolio.AvatarURL,
			GalleryURLs:  convertStringSlice(d.Portfolio.GalleryURLs),
			Introduction: &d.Portfolio.Introduction,
			SocialLinks:  convertStringSlice(d.Portfolio.SocialLinks),
		},
	}
}
