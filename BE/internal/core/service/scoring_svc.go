package service

import (
	"context"
	"errors"
	"html"
	"time"

	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
	vld "cuoc_thi_hoa_hau/internal/pkg/validation"
)

type scoringService struct {
	repo           port.ScoreRepository
	contestantRepo port.ContestantRepository
}

func NewScoringService(repo port.ScoreRepository, contestantRepo port.ContestantRepository) port.ScoreService {
	return &scoringService{
		repo:           repo,
		contestantRepo: contestantRepo,
	}
}

// GetMyScores handles both Candidates (viewing their scores) and Examiners (viewing scores they gave)
// But since the interface doesn't pass Role, we might need to rely on the caller or check user type
// Ideally, the caller (Resolver) knows the Role.
// For now, let's assume this method is primarily for the "History" view.
// We can check if the userID exists in 'contestants' -> Candidate. Else -> Examiner?
// Or better: The caller should filter.
// Let's implement robustly:
func (s *scoringService) GetMyScores(ctx context.Context, userID string) ([]*domain.Score, error) {
	// Try to see if this user is a contestant
	contestant, err := s.contestantRepo.GetByUserID(ctx, userID)
	if err == nil && contestant != nil {
		// Is a contestant -> return scores received
		return s.repo.GetListByCandidateID(ctx, contestant.ID)
	}

	// Not a contestant -> assume Examiner -> return scores given
	return s.repo.GetListByExaminerID(ctx, userID)
}

func (s *scoringService) SubmitScore(ctx context.Context, examinerID string, score *domain.Score, ip, ua string) error {
	// 1. Validation
	if err := vld.ValidateScore(score.CriteriaScores); err != nil {
		return err
	}

	if score.ContestantID == "" {
		return errors.New("thiếu ID thí sinh")
	}

	// 2. Check Contestant Existence
	contestant, err := s.contestantRepo.GetByID(ctx, score.ContestantID)
	if err != nil || contestant == nil {
		return errors.New("thí sinh không tồn tại hoặc đã bị xóa")
	}

	// 3. Sanitize Comment
	score.Comment = html.EscapeString(score.Comment)

	// 4. Bind Audit Data
	score.ExaminerID = examinerID
	score.IP = ip
	score.UserAgent = ua
	score.UpdatedAt = time.Now()

	if score.CreatedAt.IsZero() {
		score.CreatedAt = time.Now()
	}

	// 5. Calculate total score
	var total float64
	for _, v := range score.CriteriaScores {
		total += v
	}
	score.TotalScore = total

	return s.repo.Upsert(ctx, score)
}
