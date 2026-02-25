/*
Package resolver chứa toàn bộ các GraphQL Resolvers.
Tác dụng của resolver.go:
- Đóng vai trò làm Dependency Injection (DI) Container tại tầng GraphQL.
- Nó "ôm" tất cả các Service vào. Khi có truy vấn GraphQL, nó gõ cửa Service tương ứng để lấy data.
*/
package resolver

import (
	"cuoc_thi_hoa_hau/internal/types"
)

type Resolver struct {
	ContestantSvc types.ContestantService
	AuthSvc       types.AuthService
	ScoreSvc      types.ScoreService
	ScheduleSvc   types.ScheduleService
	FeedbackSvc   types.FeedbackService
	CacheSvc      types.CacheService
}
