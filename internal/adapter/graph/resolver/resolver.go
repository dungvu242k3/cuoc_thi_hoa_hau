package resolver

import (
	"cuoc_thi_hoa_hau/internal/adapter/graph/generated"
	"cuoc_thi_hoa_hau/internal/core/port"
)

type Resolver struct {
	ContestantService port.ContestantService
	ScheduleService   port.ScheduleService
	FeedbackService   port.FeedbackService
	ScoringService    port.ScoringService
}

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
