package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"

	"github.com/navikt/nada-backend/pkg/graph/models"
)

// Polly is the resolver for the polly field.
func (r *queryResolver) Polly(ctx context.Context, q string) ([]*models.QueryPolly, error) {
	return r.pollyAPI.SearchPolly(ctx, q)
}
