package resolver

import (
	"cuoc_thi_hoa_hau/internal/core/port"
)

type Resolver struct {
	ContestantSvc port.ContestantService
	AuthSvc       port.AuthService
	ScoreSvc      port.ScoreService
	ScheduleSvc   port.ScheduleService
	FeedbackSvc   port.FeedbackService
	CacheSvc      port.CacheService
}
