package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/graph/generated"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

func (r *collectionResolver) Elements(ctx context.Context, obj *models.Collection) ([]models.CollectionElement, error) {
	return r.repo.GetCollectionElements(ctx, obj.ID)
}

func (r *mutationResolver) CreateCollection(ctx context.Context, input models.NewCollection) (*models.Collection, error) {
	if err := ensureUserInGroup(ctx, input.Group); err != nil {
		return nil, err
	}
	return r.repo.CreateCollection(ctx, input)
}

func (r *mutationResolver) UpdateCollection(ctx context.Context, id uuid.UUID, input models.UpdateCollection) (*models.Collection, error) {
	col, err := r.repo.GetCollection(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureUserInGroup(ctx, col.Owner.Group); err != nil {
		return nil, err
	}
	return r.repo.UpdateCollection(ctx, id, input)
}

func (r *mutationResolver) DeleteCollection(ctx context.Context, id uuid.UUID) (bool, error) {
	col, err := r.repo.GetCollection(ctx, id)
	if err != nil {
		return false, err
	}
	if err := ensureUserInGroup(ctx, col.Owner.Group); err != nil {
		return false, err
	}
	return true, r.repo.DeleteCollection(ctx, id)
}

func (r *mutationResolver) AddToCollection(ctx context.Context, id uuid.UUID, elementID uuid.UUID, elementType models.CollectionElementType) (bool, error) {
	col, err := r.repo.GetCollection(ctx, id)
	if err != nil {
		return false, err
	}
	if err := ensureUserInGroup(ctx, col.Owner.Group); err != nil {
		return false, err
	}
	return true, r.repo.AddToCollection(ctx, id, elementID, elementType.String())
}

func (r *mutationResolver) RemoveFromCollection(ctx context.Context, id uuid.UUID, elementID uuid.UUID, elementType models.CollectionElementType) (bool, error) {
	col, err := r.repo.GetCollection(ctx, id)
	if err != nil {
		return false, err
	}
	if err := ensureUserInGroup(ctx, col.Owner.Group); err != nil {
		return false, err
	}
	return true, r.repo.RemoveFromCollection(ctx, id, elementID, elementType.String())
}

func (r *queryResolver) Collections(ctx context.Context, limit *int, offset *int) ([]*models.Collection, error) {
	l, o := pagination(limit, offset)
	return r.repo.GetCollections(ctx, l, o)
}

func (r *queryResolver) Collection(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	return r.repo.GetCollection(ctx, id)
}

// Collection returns generated.CollectionResolver implementation.
func (r *Resolver) Collection() generated.CollectionResolver { return &collectionResolver{r} }

type collectionResolver struct{ *Resolver }