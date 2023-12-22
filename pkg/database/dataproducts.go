package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/database/gensql"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

func (r *Repo) GetDataproducts(ctx context.Context, limit, offset int) ([]*models.Dataproduct, error) {
	res, err := r.querier.GetDataproducts(ctx, gensql.GetDataproductsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, fmt.Errorf("getting dataproducts from database: %w", err)
	}
	return dataproductsFromSQL(res)
}

func (r *Repo) GetDataproductsByUserAccess(ctx context.Context, user string) ([]*models.Dataproduct, error) {
	// todo: necessary?
	return nil, nil
}

func (r *Repo) GetDataproductsByGroups(ctx context.Context, groups []string) ([]*models.Dataproduct, error) {
	res, err := r.querier.GetDataproductsByGroups(ctx, groups)
	if err != nil {
		return nil, fmt.Errorf("getting dataproducts by group from database: %w", err)
	}
	return dataproductsFromSQL(res)
}

func (r *Repo) GetDataproductByProductArea(ctx context.Context, teamIDs []string) ([]*models.Dataproduct, error) {
	dps, err := r.querier.GetDataproductsByProductArea(ctx, teamIDs)
	if err != nil {
		return nil, err
	}
	return dataproductsFromSQL(dps)
}

func (r *Repo) GetDataproductByTeam(ctx context.Context, teamID string) ([]*models.Dataproduct, error) {
	dps, err := r.querier.GetDataproductsByTeam(ctx, sql.NullString{String: teamID, Valid: true})
	if err != nil {
		return nil, err
	}

	return dataproductsFromSQL(dps)
}

func (r *Repo) CreateDataproduct(ctx context.Context, dp models.NewDataproduct, user *auth.User) (*models.Dataproduct, error) {
	dataproduct, err := r.querier.CreateDataproduct(ctx, gensql.CreateDataproductParams{
		Name:                  dp.Name,
		Description:           ptrToNullString(dp.Description),
		OwnerGroup:            dp.Group,
		OwnerTeamkatalogenUrl: ptrToNullString(dp.TeamkatalogenURL),
		Slug:                  slugify(dp.Slug, dp.Name),
		TeamContact:           ptrToNullString(dp.TeamContact),
		TeamID:                ptrToNullString(dp.TeamID),
	})
	if err != nil {
		return nil, err
	}

	return minimalDataproductFromSQL(dataproduct), nil
}

func (r *Repo) UpdateDataproduct(ctx context.Context, id uuid.UUID, new models.UpdateDataproduct) error {
	_, err := r.querier.UpdateDataproduct(ctx, gensql.UpdateDataproductParams{
		Name:                  new.Name,
		Description:           ptrToNullString(new.Description),
		ID:                    id,
		OwnerTeamkatalogenUrl: ptrToNullString(new.TeamkatalogenURL),
		TeamContact:           ptrToNullString(new.TeamContact),
		Slug:                  slugify(new.Slug, new.Name),
		TeamID:                ptrToNullString(new.TeamID),
	})
	if err != nil {
		return fmt.Errorf("updating dataproduct in database: %w", err)
	}

	return nil
}

func (r *Repo) DeleteDataproduct(ctx context.Context, id uuid.UUID) error {
	if err := r.querier.DeleteDataproduct(ctx, id); err != nil {
		return fmt.Errorf("deleting dataproduct from database: %w", err)
	}

	return nil
}

func (r *Repo) GetBigqueryDatasources(ctx context.Context) ([]gensql.DatasourceBigquery, error) {
	return r.querier.GetBigqueryDatasources(ctx)
}

func (r *Repo) DataproductKeywords(ctx context.Context, prefix string) ([]*models.Keyword, error) {
	kws, err := r.querier.DataproductKeywords(ctx, prefix)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.Keyword, len(kws))
	for i, kw := range kws {
		ret[i] = &models.Keyword{
			Keyword: kw.Keyword,
			Count:   int(kw.Count),
		}
	}
	return ret, nil
}

func (r *Repo) DataproductGroupStats(ctx context.Context, limit, offset int) ([]*models.GroupStats, error) {
	stats, err := r.querier.DataproductGroupStats(ctx, gensql.DataproductGroupStatsParams{
		Lim:  int32(limit),
		Offs: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	ret := make([]*models.GroupStats, len(stats))
	for i, s := range stats {
		ret[i] = &models.GroupStats{
			Email:        s.Group,
			Dataproducts: int(s.Count),
		}
	}
	return ret, nil
}

func (r *Repo) GetDataproduct(ctx context.Context, id uuid.UUID) (*models.Dataproduct, error) {
	sqldp, err := r.querier.GetDataproduct(ctx, id)
	if err != nil {
		return nil, err
	}

	dp, err := dataproductsFromSQL(sqldp)
	if err != nil {
		return nil, err
	}
	if len(dp) == 0 {
		return nil, fmt.Errorf("GetDataproduct: no dataproduct with id %s", id)
	}
	return dp[0], nil
}

func dataproductsFromSQL(dprows []gensql.DataproductCompleteView) ([]*models.Dataproduct, error) {
	datasets, err := datasetsFromSQL(dprows)
	if err != nil {
		return nil, err
	}

	dataproducts := []*models.Dataproduct{}

	for _, dprow := range dprows {
		var dataproduct *models.Dataproduct

		for _, dp := range dataproducts {
			if dp.ID == dprow.DataproductID {
				dataproduct = dp
				break
			}
		}
		if dataproduct == nil {
			dataproduct = &models.Dataproduct{
				ID:           dprow.DataproductID,
				Name:         dprow.DpName,
				Created:      dprow.DpCreated,
				LastModified: dprow.DpLastModified,
				Description:  nullStringToPtr(dprow.DpDescription),
				Slug:         dprow.DpSlug,
				Owner: &models.Owner{
					Group:            dprow.DpGroup,
					TeamkatalogenURL: nullStringToPtr(dprow.TeamkatalogenUrl),
					TeamContact:      nullStringToPtr(dprow.TeamContact),
					TeamID:           nullStringToPtr(dprow.TeamID),
				},
			}
			dpdatasets := []*models.Dataset{}
			for _, ds := range datasets {
				if ds.DataproductID == dataproduct.ID {
					dpdatasets = append(dpdatasets, ds)
				}
			}

			keywordsMap := make(map[string]bool)
			for _, ds := range dpdatasets {
				for _, k := range ds.Keywords {
					keywordsMap[k] = true
				}
			}
			keywords := []string{}
			for k := range keywordsMap {
				keywords = append(keywords, k)
			}

			dataproduct.Datasets = dpdatasets
			dataproduct.Keywords = keywords
			dataproducts = append(dataproducts, dataproduct)
		}
	}
	return dataproducts, nil
}

func minimalDataproductFromSQL(dp gensql.Dataproduct) *models.Dataproduct {
	return &models.Dataproduct{
		ID:           dp.ID,
		Name:         dp.Name,
		Created:      dp.Created,
		LastModified: dp.LastModified,
		Description:  nullStringToPtr(dp.Description),
		Slug:         dp.Slug,
		Owner: &models.Owner{
			Group:            dp.Group,
			TeamkatalogenURL: nullStringToPtr(dp.TeamkatalogenUrl),
			TeamContact:      nullStringToPtr(dp.TeamContact),
			TeamID:           nullStringToPtr(dp.TeamID),
		},
	}
}
