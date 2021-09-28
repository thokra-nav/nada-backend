package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/navikt/datakatalogen/backend/database/gensql"
	"github.com/navikt/datakatalogen/backend/openapi"
	"net/url"

	// Pin version of sqlc and goose for cli
	_ "github.com/kyleconroy/sqlc"
	_ "github.com/pressly/goose/v3"
)

type Repo struct {
	querier gensql.Querier
}

func New(querier gensql.Querier) (*Repo, error) {
	return &Repo{
		querier: querier,
	}, nil
}

func slugify(maybeslug *string, fallback string) string {
	if maybeslug != nil {
		return *maybeslug
	}
	// TODO(thokra): Smartify this?
	return url.PathEscape(fallback)
}

func (r *Repo) CreateDataproduct(ctx context.Context, dp openapi.NewDataproduct) (*openapi.Dataproduct, error) {
	var keywords []string
	if dp.Keyword != nil {
		keywords = *dp.Keyword
	}
	res, err := r.querier.CreateDataproduct(ctx, gensql.CreateDataproductParams{
		Name:        dp.Name,
		Description: ptrToNullString(dp.Description),
		Slug:        slugify(dp.Slug, dp.Name),
		Repo:        ptrToNullString(dp.Repo),
		Team:        dp.Owner.Team,
		Keywords:    keywords,
	})
	if err != nil {
		return nil, err
	}

	return dataproductFromSQL(res), nil
}

func (r *Repo) GetDataproducts(ctx context.Context) ([]*openapi.Dataproduct, error) {
	dataproducts := []*openapi.Dataproduct{}

	res, err := r.querier.GetDataproducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting dataproducts from database: %w", err)
	}

	for _, entry := range res {
		dataproducts = append(dataproducts, dataproductFromSQL(entry))
	}

	return dataproducts, nil
}

func (r *Repo) GetDataproduct(ctx context.Context, id string) (*openapi.Dataproduct, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parsing uuid: %w", err)
	}
	res, err := r.querier.GetDataproduct(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("getting dataproduct from database: %w", err)
	}

	return dataproductFromSQL(res), nil
}

func (r *Repo) DeleteDataproduct(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parsing uuid: %w", err)
	}

	if err := r.querier.DeleteDataproduct(ctx, uid); err != nil {
		return fmt.Errorf("deleting dataproduct from database: %w", err)
	}

	return nil
}

func (r *Repo) UpdateDataproduct(ctx context.Context, id string, new openapi.NewDataproduct) (*openapi.Dataproduct, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parsing uuid: %w", err)
	}

	res, err := r.querier.UpdateDataproduct(ctx, gensql.UpdateDataproductParams{
		Name:        new.Name,
		Description: ptrToNullString(new.Description),
		Slug:        slugify(new.Slug, new.Name),
		Repo:        ptrToNullString(new.Repo),
		Team:        new.Owner.Team,
		Keywords:    *new.Keyword,
		ID:          uid,
	})

	if err != nil {
		return nil, fmt.Errorf("updating dataproduct in database: %w", err)
	}

	return dataproductFromSQL(res), nil
}

func dataproductFromSQL(dataproduct gensql.Dataproduct) *openapi.Dataproduct {
	return &openapi.Dataproduct{
		Id:           dataproduct.ID.String(),
		Name:         dataproduct.Name,
		Created:      dataproduct.Created,
		LastModified: dataproduct.LastModified,
		Description:  nullStringToPtr(dataproduct.Description),
		Keyword:      &dataproduct.Keywords,
		Owner: openapi.Owner{
			Team: dataproduct.Team,
		},
		Repo: nullStringToPtr(dataproduct.Repo),
		Slug: dataproduct.Slug,
	}
}