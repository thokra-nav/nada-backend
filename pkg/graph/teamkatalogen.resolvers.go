package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"

	"github.com/navikt/nada-backend/pkg/graph/models"
)

// Teamkatalogen is the resolver for the teamkatalogen field.
func (r *queryResolver) Teamkatalogen(ctx context.Context, q []string) ([]*models.TeamkatalogenResult, error) {
	var teamkatalogenResult []*models.TeamkatalogenResult
	if len(q) == 0 {
		q = []string{""}
	}
	for _, gcpGroup := range q {
		tr, err := r.teamkatalogen.Search(ctx, gcpGroup)
		if err != nil {
			return teamkatalogenResult, err
		}
		teamkatalogenResult = append(teamkatalogenResult, tr...)
	}
	return teamkatalogenResult, nil
}
