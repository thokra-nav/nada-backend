package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"fmt"

	"github.com/navikt/nada-backend/pkg/graph/generated"
)

// Dummy is the resolver for the dummy field.
func (r *mutationResolver) Dummy(ctx context.Context, no *string) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Version is the resolver for the version field.
func (r *queryResolver) Version(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
