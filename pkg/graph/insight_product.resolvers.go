package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"context"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

// CreateInsightProduct is the resolver for the createInsightProduct field.
func (r *mutationResolver) CreateInsightProduct(ctx context.Context, input models.NewInsightProduct) (*models.InsightProduct, error) {
	return r.repo.CreateInsightProduct(ctx, auth.GetUser(ctx).Email, input)
}

// UpdateInsightProductMetadata is the resolver for the updateInsightProductMetadata field.
func (r *mutationResolver) UpdateInsightProductMetadata(ctx context.Context, id uuid.UUID, name string, description string, typeArg string, link string, keywords []string, teamkatalogenURL *string, productAreaID *string, teamID *string, group string) (*models.InsightProduct, error) {
	existing, err := r.repo.GetInsightProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Creator = ""
	/*
		user := auth.GetUser(ctx)
		if !user.GoogleGroups.Contains(existing.Group) {
			return nil, ErrUnauthorized
		}
	*/

	insightProduct, err := r.repo.UpdateInsightProductMetadata(ctx, id, name, description, keywords, teamkatalogenURL,
		productAreaID, teamID, typeArg, link)
	if err != nil {
		return nil, err
	}

	return insightProduct, nil
}

// DeleteInsightProduct is the resolver for the deleteInsightProduct field.
func (r *mutationResolver) DeleteInsightProduct(ctx context.Context, id uuid.UUID) (bool, error) {
	s, err := r.repo.GetInsightProduct(ctx, id)
	if err != nil {
		return false, err
	}

	s.Creator = ""
	/*
		user := auth.GetUser(ctx)
		if !user.GoogleGroups.Contains(s.Group) {
			return false, ErrUnauthorized
		}
	*/
	if err = r.repo.DeleteInsightProduct(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}

// InsightProduct is the resolver for the InsightProduct field.
func (r *queryResolver) InsightProduct(ctx context.Context, id uuid.UUID) (*models.InsightProduct, error) {
	return r.repo.GetInsightProduct(ctx, id)
}