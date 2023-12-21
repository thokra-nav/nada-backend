package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/database/gensql"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

func (r *Repo) GetDataproducts(ctx context.Context, limit, offset int) ([]*models.Dataproduct, error) {
	dataproducts := []*models.Dataproduct{}

	res, err := r.querier.GetDataproducts(ctx, gensql.GetDataproductsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, fmt.Errorf("getting dataproducts from database: %w", err)
	}

	for _, entry := range res {
		dataproducts = append(dataproducts, dataproductFromSQL(entry))
	}

	return dataproducts, nil
}

func (r *Repo) GetDataproductsByUserAccess(ctx context.Context, user string) ([]*models.Dataproduct, error) {
	// todo: necessary?
	return nil, nil
}

func (r *Repo) GetDataproductsByGroups(ctx context.Context, groups []string) ([]*models.Dataproduct, error) {
	dps := []*models.Dataproduct{}

	res, err := r.querier.GetDataproductsByGroups(ctx, groups)
	if err != nil {
		return nil, fmt.Errorf("getting dataproducts by group from database: %w", err)
	}

	for _, entry := range res {
		dps = append(dps, dataproductFromSQL(entry))
	}

	return dps, nil
}

func (r *Repo) GetDataproductByProductArea(ctx context.Context, teamIDs []string) ([]*models.Dataproduct, error) {
	dps, err := r.querier.GetDataproductsByProductArea(ctx, teamIDs)
	if err != nil {
		return nil, err
	}

	dpsGraph := make([]*models.Dataproduct, len(dps))
	for idx, dp := range dps {
		dpsGraph[idx] = dataproductFromSQL(dp)
	}

	return dpsGraph, nil
}

func (r *Repo) GetDataproductByTeam(ctx context.Context, teamID string) ([]*models.Dataproduct, error) {
	dps, err := r.querier.GetDataproductsByTeam(ctx, sql.NullString{String: teamID, Valid: true})
	if err != nil {
		return nil, err
	}

	dpsGraph := make([]*models.Dataproduct, len(dps))
	for idx, dp := range dps {
		dpsGraph[idx] = dataproductFromSQL(dp)
	}

	return dpsGraph, nil
}

func (r *Repo) GetDataproduct(ctx context.Context, id uuid.UUID) (*models.Dataproduct, error) {
	res, err := r.querier.GetDataproduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting dataproduct from database: %w", err)
	}

	return dataproductFromSQL(res), nil
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

	return dataproductFromSQL(dataproduct), nil
}

func (r *Repo) UpdateDataproduct(ctx context.Context, id uuid.UUID, new models.UpdateDataproduct) (*models.Dataproduct, error) {
	res, err := r.querier.UpdateDataproduct(ctx, gensql.UpdateDataproductParams{
		Name:                  new.Name,
		Description:           ptrToNullString(new.Description),
		ID:                    id,
		OwnerTeamkatalogenUrl: ptrToNullString(new.TeamkatalogenURL),
		TeamContact:           ptrToNullString(new.TeamContact),
		Slug:                  slugify(new.Slug, new.Name),
		TeamID:                ptrToNullString(new.TeamID),
	})
	if err != nil {
		return nil, fmt.Errorf("updating dataproduct in database: %w", err)
	}

	return dataproductFromSQL(res), nil
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

func (r *Repo) GetDataproductComplete(ctx context.Context, id uuid.UUID) (*models.DataproductComplete, error) {
	sqldp, err := r.querier.GetDataproductComplete(ctx, id)
	if err != nil {
		return nil, err
	}

	return dataproductCompleteFromSQL(sqldp)
}

func dataproductFromSQL(dp gensql.Dataproduct) *models.Dataproduct {
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

func dataproductCompleteFromSQL(dprows []gensql.GetDataproductCompleteRow) (*models.DataproductComplete, error) {
	dp := models.Dataproduct{
		ID:           dprows[0].ID,
		Name:         dprows[0].Name,
		Created:      dprows[0].Created,
		LastModified: dprows[0].LastModified,
		Description:  nullStringToPtr(dprows[0].Description),
		Slug:         dprows[0].Slug,
		Owner: &models.Owner{
			Group:            dprows[0].Group,
			TeamkatalogenURL: nullStringToPtr(dprows[0].TeamkatalogenUrl),
			TeamContact:      nullStringToPtr(dprows[0].TeamContact),
			TeamID:           nullStringToPtr(dprows[0].TeamID),
		},
	}

	datasets := []*models.DatasetComplete{}
	for _, dprow := range dprows {
		if dprow.DsID.UUID == uuid.Nil {
			continue
		}

		piiTags := "{}"
		if dprow.PiiTags.RawMessage != nil {
			piiTags = string(dprow.PiiTags.RawMessage)
		}

		var ds *models.DatasetComplete

		for _, dsIn := range datasets {
			if dsIn.ID == dprow.DsID.UUID {
				ds = dsIn
				fmt.Println("found")
				break
			}
		}
		if ds == nil {
			ds = &models.DatasetComplete{
				Dataset: models.Dataset{
					ID:            dprow.DsID.UUID,
					Name:          dprow.DsName.String,
					Created:       dprow.DsrcCreated,
					LastModified:  dprow.DsrcLastModified,
					Description:   nullStringToPtr(dprow.DsDescription),
					Slug:          dprow.DsSlug.String,
					Keywords:      dprow.Keywords,
					DataproductID: dp.ID,
				},
				Owner:    dp.Owner,
				Mappings: []models.MappingService{},
				Access:   []*models.Access{},
				Services: &models.DatasetServices{},
			}
			datasets = append(datasets, ds)
		}

		if dprow.DsrcID != uuid.Nil {
			var schema []*models.TableColumn
			if dprow.DsrcSchema.Valid {
				if err := json.Unmarshal(dprow.DsrcSchema.RawMessage, &schema); err != nil {
					return nil, fmt.Errorf("unmarshalling schema: %w", err)
				}
			}

			dsrc := models.BigQueryComplete{
				BigQuery: models.BigQuery{
					ID:            dprow.DsrcID,
					DatasetID:     dprow.DsID.UUID,
					ProjectID:     dprow.ProjectID,
					Dataset:       dprow.Dataset,
					Table:         dprow.TableName,
					TableType:     models.BigQueryType(dprow.TableType),
					Created:       dprow.DsrcCreated,
					LastModified:  dprow.DsrcLastModified,
					Expires:       nullTimeToPtr(dprow.DsrcExpires),
					Description:   dprow.DsrcDescription.String,
					PiiTags:       &piiTags,
					MissingSince:  nullTimeToPtr(dprow.DsrcMissingSince),
					PseudoColumns: dprow.PseudoColumns,
				},
				Schema: schema,
			}
			ds.Datasource = dsrc
		}

		if len(dprow.Services) > 0 {
			for _, service := range dprow.Services {
				exist := false
				for _, mapping := range ds.Mappings {
					if mapping.String() == service {
						exist = true
						break
					}
				}
				if !exist {
					ds.Mappings = append(ds.Mappings, models.MappingService(service))
				}
			}
		}

		if dprow.DaID.Valid {
			exist := false
			for _, dsAccess := range ds.Access {
				if dsAccess.ID == dprow.DaID.UUID {
					exist = true
					break
				}
			}
			if !exist {
				access := &models.Access{
					ID:              dprow.DaID.UUID,
					Subject:         dprow.DaSubject.String,
					Granter:         dprow.DaGranter.String,
					Expires:         nullTimeToPtr(dprow.DaExpires),
					Created:         dprow.DaCreated.Time,
					Revoked:         nullTimeToPtr(dprow.DaRevoked),
					DatasetID:       dprow.DsID.UUID,
					AccessRequestID: nullUUIDToUUIDPtr(dprow.AccessRequestID),
				}
				ds.Access = append(ds.Access, access)
			}
		}

		if ds.Services == nil && dprow.MmDatabaseID.Valid {
			svc := &models.DatasetServices{}
			base := "https://metabase.intern.dev.nav.no/browse/%v"
			if os.Getenv("NAIS_CLUSTER_NAME") == "prod-gcp" {
				base = "https://metabase.intern.nav.no/browse/%v"
			}
			url := fmt.Sprintf(base, dprow.MmDatabaseID.Int32)
			svc.Metabase = &url
			ds.Services = svc
		}
	}
	keywordsMap := make(map[string]bool)
	for _, ds := range datasets {
		for _, k := range ds.Keywords {
			keywordsMap[k] = true
		}
	}
	keywords := []string{}
	for k := range keywordsMap {
		keywords = append(keywords, k)
	}

	dpcomplete := &models.DataproductComplete{
		Dataproduct: dp,
		Datasets:    datasets,
		Keywords:    keywords,
	}
	return dpcomplete, nil
}
