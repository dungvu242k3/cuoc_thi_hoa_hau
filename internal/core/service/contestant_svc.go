package service

import (
	"context"
	"errors"
	"html"
	"log"
	"time"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
)

type contestantService struct {
	repo port.ContestantRepository
}

func NewContestantService(repo port.ContestantRepository) port.ContestantService {
	return &contestantService{repo: repo}
}

func (s *contestantService) CreateProfile(ctx context.Context, userID string, input *domain.Contestant) (*domain.Contestant, error) {
	//  Check Trùng
	exist, _ := s.repo.GetByUserID(ctx, userID)
	if exist != nil {
		return nil, errors.New("bạn đã có hồ sơ, không thể tạo thêm")
	}
	if ok, _ := s.repo.CheckIdentifyCard(ctx, input.IdentifyCard); ok {
		return nil, errors.New("số CCCD này đã được sử dụng")
	}

	// Check Tuổi Chiều cao
	age := time.Now().Year() - input.PersonalInfo.DateOfBirth.Year()
	if age < 18 {
		return nil, errors.New("thí sinh phải đủ 18 tuổi")
	}
	if input.PhysicalInfo.Height < 160 {
		return nil, errors.New("chiều cao phải trên 1m60")
	}

	// Setup Default
	input.UserID = userID
	input.Status = domain.StatusDraft
	input.IsPublic = false
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	input.Portfolio.Introduction = html.EscapeString(input.Portfolio.Introduction)

	if err := s.repo.Create(ctx, input); err != nil {
		return nil, err
	}
	return input, nil
}

func (s *contestantService) UpdateProfile(ctx context.Context, userID string, input *domain.Contestant) (*domain.Contestant, error) {
	current, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("hồ sơ không tồn tại")
	}

	// CHẶN SỬA NẾU ĐÃ NỘP/DUYỆT
	if current.Status == domain.StatusPending || current.Status == domain.StatusApproved {
		return nil, errors.New("hồ sơ đang chờ duyệt hoặc đã duyệt, không được phép chỉnh sửa")
	}

	// Personal Info
	if input.PersonalInfo.FullName != "" {
		current.PersonalInfo.FullName = input.PersonalInfo.FullName
	}
	if !input.PersonalInfo.DateOfBirth.IsZero() {
		current.PersonalInfo.DateOfBirth = input.PersonalInfo.DateOfBirth
	}
	if input.PersonalInfo.Nationality != "" {
		current.PersonalInfo.Nationality = input.PersonalInfo.Nationality
	}
	if input.PersonalInfo.Phone != "" {
		current.PersonalInfo.Phone = input.PersonalInfo.Phone
	}
	if input.PersonalInfo.Email != "" {
		current.PersonalInfo.Email = input.PersonalInfo.Email
	}
	if input.PersonalInfo.Address != "" {
		current.PersonalInfo.Address = input.PersonalInfo.Address
	}

	// Physical Info
	if input.PhysicalInfo.Height > 0 {
		current.PhysicalInfo.Height = input.PhysicalInfo.Height
	}
	if input.PhysicalInfo.Weight > 0 {
		current.PhysicalInfo.Weight = input.PhysicalInfo.Weight
	}
	if input.PhysicalInfo.Measurements != "" {
		current.PhysicalInfo.Measurements = input.PhysicalInfo.Measurements
	}

	// Skill & Edu
	if input.SkillEdu.EducationLevel != "" {
		current.SkillEdu.EducationLevel = input.SkillEdu.EducationLevel
	}
	if len(input.SkillEdu.Languages) > 0 {
		current.SkillEdu.Languages = input.SkillEdu.Languages
	}
	if len(input.SkillEdu.Skills) > 0 {
		current.SkillEdu.Skills = input.SkillEdu.Skills
	}

	// Portfolio
	if input.Portfolio.AvatarURL != "" {
		current.Portfolio.AvatarURL = input.Portfolio.AvatarURL
	}
	if len(input.Portfolio.GalleryURLs) > 0 {
		current.Portfolio.GalleryURLs = input.Portfolio.GalleryURLs
	}
	if input.Portfolio.Introduction != "" {
		current.Portfolio.Introduction = html.EscapeString(input.Portfolio.Introduction)
	}
	if len(input.Portfolio.SocialLinks) > 0 {
		current.Portfolio.SocialLinks = input.Portfolio.SocialLinks
	}

	current.UpdatedAt = time.Now()

	// Log Action (Secure: UserID only, no PII)
	log.Printf("[AUDIT] User %s updated profile at %s", userID, time.Now().Format(time.RFC3339))

	return current, s.repo.Update(ctx, current)
}

func (s *contestantService) SubmitProfile(ctx context.Context, userID string) error {
	current, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Validate dữ liệu bắt buộc trước khi nộp
	if current.Portfolio.AvatarURL == "" {
		return errors.New("vui lòng cập nhật ảnh đại diện")
	}

	current.Status = domain.StatusPending
	current.UpdatedAt = time.Now()
	return s.repo.Update(ctx, current)
}

func (s *contestantService) GetMyProfile(ctx context.Context, userID string) (*domain.Contestant, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *contestantService) GetPublicList(ctx context.Context, limit, offset int64) ([]*domain.Contestant, int64, error) {
	return s.repo.GetPublicList(ctx, limit, offset)
}

func (s *contestantService) GetPublicDetail(ctx context.Context, id string) (*domain.Contestant, error) {
	return s.repo.GetPublicDetail(ctx, id)
}
