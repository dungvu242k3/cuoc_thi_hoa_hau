/*
Package service chứa trọn vẹn "Nghiệp vụ cốt lõi" (Business Logic).
Tác dụng của contestant.go:
- Không được chứa lệnh SQL/BSON (đã giao cho DAO).
- Không được biết nó đang chạy qua HTTP hay gRPC (đã giao cho Handler/Resolver).
- Chỉ tập trung xử lý: Điều kiện tuổi > 18, kiểm tra trùng CCCD, duyệt hồ sơ.
*/
package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"time"

	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"
)

type contestantService struct {
	repo types.ContestantRepository
}

func NewContestantService(repo types.ContestantRepository) types.ContestantService {
	return &contestantService{repo: repo}
}

// CreateProfile là nghiệp vụ Nộp hồ sơ.
// Nó gọi qua DAO để kiểm tra xem User này đã nộp bao giờ chưa, CCCD có bị trùng ai không.
// Nếu thỏa mãn, nó mới tạo ra bản ghi ContestantStatusPending chờ Admin duyệt.
func (s *contestantService) CreateProfile(ctx context.Context, userID string, input *model.Contestant) (*model.Contestant, error) {
	exist, _ := s.repo.GetByUserID(ctx, userID)
	if exist != nil {
		return nil, errors.New("bạn đã có hồ sơ, không thể tạo thêm")
	}
	if ok, _ := s.repo.CheckIdentifyCard(ctx, input.PersonalInfo.IdentityCard); ok {
		return nil, errors.New("số CCCD này đã được sử dụng")
	}

	const (
		MinAge    = 18
		MinHeight = 160.0
	)

	age := time.Now().Year() - input.PersonalInfo.DateOfBirth.Year()
	if age < MinAge {
		return nil, fmt.Errorf("thí sinh phải đủ %d tuổi", MinAge)
	}
	if input.PhysicalInfo.Height < MinHeight {
		return nil, fmt.Errorf("chiều cao phải trên %.0fcm", MinHeight)
	}

	count, _ := s.repo.Count(ctx)
	input.SBD = fmt.Sprintf("%03d", count+1)

	input.UserID = userID
	input.Status = model.ContestantStatusPending
	input.IsPublic = false
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	input.Portfolio.Introduction = html.EscapeString(input.Portfolio.Introduction)

	if err := s.repo.Create(ctx, input); err != nil {
		return nil, err
	}
	return input, nil
}

// UpdateProfile cho phép thí sinh sửa hồ sơ NẾU hồ sơ chưa được duyệt (hồ sơ đang Draft/Pending).
// Hàm này là một ví dụ về Patch Update (chỉ đè dữ liệu mới lên các trường khác Rỗng).
func (s *contestantService) UpdateProfile(ctx context.Context, userID string, input *model.Contestant) (*model.Contestant, error) {
	current, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("lỗi truy vấn hồ sơ")
	}
	if current == nil {
		return nil, errors.New("hồ sơ không tồn tại, vui lòng tạo mới")
	}

	if current.Status == model.ContestantStatusApproved || current.Status == model.ContestantStatusRejected {
		return nil, fmt.Errorf("hồ sơ đang ở trạng thái '%s' (đã xử lý), không được phép chỉnh sửa", current.Status)
	}

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
	if input.PersonalInfo.Job != "" {
		current.PersonalInfo.Job = input.PersonalInfo.Job
	}

	if input.PhysicalInfo.Height > 0 {
		current.PhysicalInfo.Height = input.PhysicalInfo.Height
	}
	if input.PhysicalInfo.Weight > 0 {
		current.PhysicalInfo.Weight = input.PhysicalInfo.Weight
	}
	if input.PhysicalInfo.Measurements != "" {
		current.PhysicalInfo.Measurements = input.PhysicalInfo.Measurements
	}

	if input.SkillEdu.EducationLevel != "" {
		current.SkillEdu.EducationLevel = input.SkillEdu.EducationLevel
	}
	if len(input.SkillEdu.Languages) > 0 {
		current.SkillEdu.Languages = input.SkillEdu.Languages
	}
	if len(input.SkillEdu.Skills) > 0 {
		current.SkillEdu.Skills = input.SkillEdu.Skills
	}

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

	log.Printf("[AUDIT] User %s updated profile at %s", userID, time.Now().Format(time.RFC3339))

	return current, s.repo.Update(ctx, current)
}

func (s *contestantService) SubmitProfile(ctx context.Context, userID string) error {
	current, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if current.Portfolio.AvatarURL == "" {
		return errors.New("vui lòng cập nhật ảnh đại diện")
	}

	current.Status = model.ContestantStatusPending
	current.UpdatedAt = time.Now()
	return s.repo.Update(ctx, current)
}

func (s *contestantService) GetMyProfile(ctx context.Context, userID string) (*model.Contestant, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *contestantService) GetPublicList(ctx context.Context, limit, offset int64) ([]*model.Contestant, int64, error) {
	return s.repo.GetPublicList(ctx, limit, offset)
}

func (s *contestantService) GetPublicDetail(ctx context.Context, id string) (*model.Contestant, error) {
	return s.repo.GetPublicDetail(ctx, id)
}

func (s *contestantService) DeleteProfile(ctx context.Context, userID string) error {
	current, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("hồ sơ không tồn tại")
	}

	current.Status = model.ContestantStatusDeleted
	current.UpdatedAt = time.Now()

	return s.repo.Update(ctx, current)
}

// ApproveContestant (Duyệt hồ sơ) là đặc quyền của Ban Tổ Chức (Admin).
// Nếu isApproved là true -> Đánh dấu Approved và hiển thị ra Public (IsPublic = true).
// Nếu false -> Rejected.
func (s *contestantService) ApproveContestant(ctx context.Context, sub string, id string, isApproved bool) (*model.Contestant, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("hồ sơ không tồn tại")
	}

	if isApproved {
		current.Status = model.ContestantStatusApproved
		current.IsPublic = true
	} else {
		current.Status = model.ContestantStatusRejected
		current.IsPublic = false
	}
	current.UpdatedAt = time.Now()

	log.Printf("[AUDIT] User %s change status of contestant %s to %v", sub, id, isApproved)

	if err := s.repo.Update(ctx, current); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *contestantService) GetAdminList(ctx context.Context, limit int64, offset int64, status *string) ([]*model.Contestant, int64, error) {
	filter := make(map[string]interface{})

	if status != nil {
		filter["status"] = *status
	}

	return s.repo.GetList(ctx, limit, offset, filter)
}

func (s *contestantService) VoteForContestant(ctx context.Context, userID, id, ip, userAgent string) error {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("thí sinh không tồn tại")
	}
	if current.Status != model.ContestantStatusApproved {
		return errors.New("chỉ có thể bình chọn cho thí sinh đã qua vòng duyệt hồ sơ")
	}

	limitReached, err := s.repo.CheckIPLimit(ctx, ip, id)
	if err != nil {
		return err
	}
	if limitReached {
		return errors.New("địa chỉ IP này đã bình chọn cho thí sinh này rồi")
	}

	hasVoted, err := s.repo.HasVoted(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasVoted {
		return errors.New("bạn đã bình chọn cho thí sinh này rồi")
	}

	if err := s.repo.RecordVote(ctx, userID, id, ip, userAgent); err != nil {
		return err
	}

	return s.repo.IncrementVote(ctx, id)
}
