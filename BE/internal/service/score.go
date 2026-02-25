/*
Package service đảm nhận Business Logic.
Tác dụng của score.go:
- Cung cấp hàm SubmitScore để giám khảo nộp điểm cho thí sinh.
- Kiểm tra chặt chẽ thí sinh đó có tồn tại thật hay không trước khi ghi điểm vào DB.
- Tính tự động điểm TotalScore trước khi lưu xuống Repo.
*/
package service

import (
	"context"
	"errors"
	"html"
	"time"

	vld "cuoc_thi_hoa_hau/internal/handler/validation"
	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"
)

type scoringService struct {
	repo           types.ScoreRepository
	contestantRepo types.ContestantRepository
}

func NewScoringService(repo types.ScoreRepository, contestantRepo types.ContestantRepository) types.ScoreService {
	return &scoringService{
		repo:           repo,
		contestantRepo: contestantRepo,
	}
}

func (s *scoringService) GetMyScores(ctx context.Context, userID string) ([]*model.Score, error) {
	contestant, err := s.contestantRepo.GetByUserID(ctx, userID)
	if err == nil && contestant != nil {
		return s.repo.GetListByCandidateID(ctx, contestant.ID)
	}

	return s.repo.GetListByExaminerID(ctx, userID)
}

// SubmitScore nhận điểm do Ban Giám Khảo chấm.
// Nó thực hiện validate dữ liệu, tự động cộng tổng (TotalScore), ghi nhận IP, thiết bị (UserAgent).
// Chống giả mạo điểm số nếu chẳng may Giám khảo bị hack tài khoản.
func (s *scoringService) SubmitScore(ctx context.Context, examinerID string, score *model.Score, ip, ua string) error {
	if err := vld.ValidateScore(score.CriteriaScores); err != nil {
		return err
	}

	if score.ContestantID == "" {
		return errors.New("thiếu ID thí sinh")
	}

	contestant, err := s.contestantRepo.GetByID(ctx, score.ContestantID)
	if err != nil || contestant == nil {
		return errors.New("thí sinh không tồn tại hoặc đã bị xóa")
	}

	score.Comment = html.EscapeString(score.Comment)

	score.ExaminerID = examinerID
	score.IP = ip
	score.UserAgent = ua
	score.UpdatedAt = time.Now()

	if score.CreatedAt.IsZero() {
		score.CreatedAt = time.Now()
	}

	var total float64
	for _, v := range score.CriteriaScores {
		total += v
	}
	score.TotalScore = total

	return s.repo.Upsert(ctx, score)
}
