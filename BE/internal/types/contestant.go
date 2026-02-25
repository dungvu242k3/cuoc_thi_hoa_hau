/*
Package types chứa toàn bộ Interface để giao tiếp giữa các Layer.
Tác dụng của contestant.go:
- ContestantRepository: Hợp đồng cho tầng DAO (quy định các lệnh nhúng xuống CSDL).
- ContestantService: Hợp đồng cho tầng Service (quy định các API nghiệp vụ cấp cao).
- Tách bạch rõ 2 khái niệm Truy vấn (Repo) và Nghiệp vụ (Service).
*/
package types

import (
	"context"
	"cuoc_thi_hoa_hau/internal/model"
)

// ContestantRepository quy định các hàm thao tác DB CRUD cho đối tượng Contestant.
type ContestantRepository interface {
	Create(ctx context.Context, c *model.Contestant) error
	Update(ctx context.Context, c *model.Contestant) error
	GetByUserID(ctx context.Context, userID string) (*model.Contestant, error)
	GetByID(ctx context.Context, id string) (*model.Contestant, error)
	GetPublicList(ctx context.Context, limit int64, offset int64) ([]*model.Contestant, int64, error)
	GetPublicDetail(ctx context.Context, id string) (*model.Contestant, error)
	CheckIdentifyCard(ctx context.Context, cardID string) (bool, error)
	Count(ctx context.Context) (int64, error)
	GetList(ctx context.Context, limit, offset int64, filter map[string]interface{}) ([]*model.Contestant, int64, error)
	IncrementVote(ctx context.Context, id string) error
	HasVoted(ctx context.Context, userID, contestantID string) (bool, error)
	RecordVote(ctx context.Context, userID, contestantID, ip, userAgent string) error
	CheckIPLimit(ctx context.Context, ip string, contestantID string) (bool, error)
}

// ContestantService quy định các tác vụ tính toán, kinh doanh, tổng hợp đối với Thí sinh.
type ContestantService interface {
	CreateProfile(ctx context.Context, userID string, input *model.Contestant) (*model.Contestant, error)
	UpdateProfile(ctx context.Context, userID string, input *model.Contestant) (*model.Contestant, error)
	SubmitProfile(ctx context.Context, userID string) error
	GetMyProfile(ctx context.Context, userID string) (*model.Contestant, error)
	DeleteProfile(ctx context.Context, userID string) error
	GetPublicList(ctx context.Context, limit int64, offset int64) ([]*model.Contestant, int64, error)
	GetPublicDetail(ctx context.Context, id string) (*model.Contestant, error)
	ApproveContestant(ctx context.Context, sub string, id string, isApproved bool) (*model.Contestant, error)
	GetAdminList(ctx context.Context, limit int64, offset int64, status *string) ([]*model.Contestant, int64, error)
	VoteForContestant(ctx context.Context, userID, id, ip, userAgent string) error
}
