package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/graph/generated"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

func (r *dataproductResolver) Datasource(ctx context.Context, obj *models.Dataproduct) (models.Datasource, error) {
	return r.repo.GetBigqueryDatasource(ctx, obj.ID)
}

func (r *mutationResolver) CreateDataproduct(ctx context.Context, input models.NewDataproduct) (*models.Dataproduct, error) {
	if err := ensureUserInGroup(ctx, input.Group); err != nil {
		return nil, err
	}
	return r.repo.CreateDataproduct(ctx, input)
}

func (r *mutationResolver) UpdateDataproduct(ctx context.Context, id uuid.UUID, input models.UpdateDataproduct) (*models.Dataproduct, error) {
	dp, err := r.repo.GetDataproduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureUserInGroup(ctx, dp.Owner.Group); err != nil {
		return nil, err
	}
	return r.repo.UpdateDataproduct(ctx, id, input)
}

func (r *mutationResolver) DeleteDataproduct(ctx context.Context, id uuid.UUID) (bool, error) {
	dp, err := r.repo.GetDataproduct(ctx, id)
	if err != nil {
		return false, err
	}
	if err := ensureUserInGroup(ctx, dp.Owner.Group); err != nil {
		return false, err
	}

	return true, r.repo.DeleteDataproduct(ctx, dp.ID)
}

func (r *queryResolver) Dataproduct(ctx context.Context, id uuid.UUID) (*models.Dataproduct, error) {
	return r.repo.GetDataproduct(ctx, id)
}

func (r *queryResolver) Dataproducts(ctx context.Context, limit *int, offset *int) ([]*models.Dataproduct, error) {
	l, o := pagination(limit, offset)
	return r.repo.GetDataproducts(ctx, l, o)
}

// Dataproduct returns generated.DataproductResolver implementation.
func (r *Resolver) Dataproduct() generated.DataproductResolver { return &dataproductResolver{r} }

type dataproductResolver struct{ *Resolver }