package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/database/gensql"
	"github.com/navikt/nada-backend/pkg/graph/models"
	"github.com/sqlc-dev/pqtype"
)

func (r *Repo) GetDataset(ctx context.Context, id uuid.UUID) (*models.Dataset, error) {
	res, err := r.querier.GetDataset(ctx, uuid.NullUUID{UUID: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("getting dataset from database: %w", err)
	}

	datasets, err := datasetsFromSQL(res)
	if err != nil {
		return nil, err
	}
	if len(datasets) == 0 {
		return nil, fmt.Errorf("GetDataset: no dataset found")
	}
	return datasets[0], nil
}

func (r *Repo) CreateDataset(ctx context.Context, ds models.NewDataset, referenceDatasource *models.NewBigQuery, user *auth.User) (*models.Dataset, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	if ds.Keywords == nil {
		ds.Keywords = []string{}
	}

	querier := r.querier.WithTx(tx)
	created, err := querier.CreateDataset(ctx, gensql.CreateDatasetParams{
		Name:                     ds.Name,
		DataproductID:            ds.DataproductID,
		Description:              ptrToNullString(ds.Description),
		Pii:                      gensql.PiiLevel(ds.Pii.String()),
		Type:                     "bigquery",
		Slug:                     slugify(ds.Slug, ds.Name),
		Repo:                     ptrToNullString(ds.Repo),
		Keywords:                 ds.Keywords,
		AnonymisationDescription: ptrToNullString(ds.AnonymisationDescription),
		TargetUser:               ptrToNullString(ds.TargetUser),
	})
	if err != nil {
		return nil, err
	}

	schemaJSON, err := json.Marshal(ds.Metadata.Schema.Columns)
	if err != nil {
		return nil, fmt.Errorf("marshalling schema: %w", err)
	}

	if ds.BigQuery.PiiTags != nil && !json.Valid([]byte(*ds.BigQuery.PiiTags)) {
		return nil, fmt.Errorf("invalid pii tags, must be json map or null: %w", err)
	}

	_, err = querier.CreateBigqueryDatasource(ctx, gensql.CreateBigqueryDatasourceParams{
		DatasetID:    created.ID,
		ProjectID:    ds.BigQuery.ProjectID,
		Dataset:      ds.BigQuery.Dataset,
		TableName:    ds.BigQuery.Table,
		Schema:       pqtype.NullRawMessage{RawMessage: schemaJSON, Valid: len(schemaJSON) > 4},
		LastModified: ds.Metadata.LastModified,
		Created:      ds.Metadata.Created,
		Expires:      sql.NullTime{Time: ds.Metadata.Expires, Valid: !ds.Metadata.Expires.IsZero()},
		TableType:    string(ds.Metadata.TableType),
		PiiTags: pqtype.NullRawMessage{
			RawMessage: json.RawMessage([]byte(ptrToString(ds.BigQuery.PiiTags))),
			Valid:      len(ptrToString(ds.BigQuery.PiiTags)) > 4,
		},
		PseudoColumns: ds.PseudoColumns,
		IsReference:   false,
	})

	if err != nil {
		if err := tx.Rollback(); err != nil {
			r.log.WithError(err).Error("Rolling back dataset and datasource_bigquery transaction")
		}
		return nil, err
	}

	if len(ds.PseudoColumns) > 0 && referenceDatasource != nil {
		_, err = querier.CreateBigqueryDatasource(ctx, gensql.CreateBigqueryDatasourceParams{
			DatasetID:    created.ID,
			ProjectID:    referenceDatasource.ProjectID,
			Dataset:      referenceDatasource.Dataset,
			TableName:    referenceDatasource.Table,
			Schema:       pqtype.NullRawMessage{RawMessage: schemaJSON, Valid: len(schemaJSON) > 4},
			LastModified: ds.Metadata.LastModified,
			Created:      ds.Metadata.Created,
			Expires:      sql.NullTime{Time: ds.Metadata.Expires, Valid: !ds.Metadata.Expires.IsZero()},
			TableType:    string(ds.Metadata.TableType),
			PiiTags: pqtype.NullRawMessage{
				RawMessage: json.RawMessage([]byte(ptrToString(ds.BigQuery.PiiTags))),
				Valid:      len(ptrToString(ds.BigQuery.PiiTags)) > 4,
			},
			PseudoColumns: ds.PseudoColumns,
			IsReference:   true,
		})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				r.log.WithError(err).Error("Rolling back dataset and datasource_bigquery transaction")
			}
			return nil, err
		}
	}

	if ds.GrantAllUsers != nil && *ds.GrantAllUsers {
		_, err = querier.GrantAccessToDataset(ctx, gensql.GrantAccessToDatasetParams{
			DatasetID: created.ID,
			Expires:   sql.NullTime{},
			Subject:   emailOfSubjectToLower("group:all-users@nav.no"),
			Granter:   user.Email,
		})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				r.log.WithError(err).Error("Rolling back dataset and datasource_bigquery transaction")
			}
			return nil, err
		}
	}

	for _, keyword := range ds.Keywords {
		err = querier.CreateTagIfNotExist(ctx, keyword)
		if err != nil {
			r.log.WithError(err).Warn("failed to create tag when creating dataset in database")
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	ret := minimalDatasetFromSQL(created)
	return ret, nil
}

func (r *Repo) CreateJoinableViews(ctx context.Context, name, owner string, expires *time.Time, datasourceIDs []uuid.UUID) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	jv, err := r.querier.CreateJoinableViews(ctx, gensql.CreateJoinableViewsParams{
		Name:    name,
		Owner:   owner,
		Created: time.Now(),
		Expires: ptrToNullTime(expires),
	})
	if err != nil {
		return "", err
	}
	for _, bqid := range datasourceIDs {
		if err != nil {
			return "", err
		}

		_, err = r.querier.CreateJoinableViewsDatasource(ctx, gensql.CreateJoinableViewsDatasourceParams{
			JoinableViewID: jv.ID,
			DatasourceID:   bqid,
		})

		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return jv.ID.String(), nil
}

func (r *Repo) UpdateDataset(ctx context.Context, id uuid.UUID, new models.UpdateDataset) (*models.Dataset, error) {
	if new.Keywords == nil {
		new.Keywords = []string{}
	}

	res, err := r.querier.UpdateDataset(ctx, gensql.UpdateDatasetParams{
		Name:                     new.Name,
		Description:              ptrToNullString(new.Description),
		ID:                       id,
		Pii:                      gensql.PiiLevel(new.Pii.String()),
		Slug:                     slugify(new.Slug, new.Name),
		Repo:                     ptrToNullString(new.Repo),
		Keywords:                 new.Keywords,
		DataproductID:            *new.DataproductID,
		AnonymisationDescription: ptrToNullString(new.AnonymisationDescription),
		TargetUser:               ptrToNullString(new.TargetUser),
	})
	if err != nil {
		return nil, fmt.Errorf("updating dataset in database: %w", err)
	}

	for _, keyword := range new.Keywords {
		err = r.querier.CreateTagIfNotExist(ctx, keyword)
		if err != nil {
			r.log.WithError(err).Warn("failed to create tag when updating dataset in database")
		}
	}

	if new.PiiTags != nil && !json.Valid([]byte(*new.PiiTags)) {
		return nil, fmt.Errorf("invalid pii tags, must be json map or null: %w", err)
	}

	err = r.querier.UpdateBigqueryDatasource(ctx, gensql.UpdateBigqueryDatasourceParams{
		DatasetID: id,
		PiiTags: pqtype.NullRawMessage{
			RawMessage: json.RawMessage(ptrToString(new.PiiTags)),
			Valid:      len(ptrToString(new.PiiTags)) > 4,
		},
		PseudoColumns: new.PseudoColumns,
	})
	if err != nil {
		return nil, err
	}

	return minimalDatasetFromSQL(res), nil
}

func (r *Repo) GetBigqueryDatasource(ctx context.Context, datasetID uuid.UUID, isReference bool) (models.BigQuery, error) {
	fmt.Println("database getbigquerydatasource")
	bq, err := r.querier.GetBigqueryDatasource(ctx, gensql.GetBigqueryDatasourceParams{
		DatasetID:   datasetID,
		IsReference: isReference,
	})
	if err != nil {
		return models.BigQuery{}, err
	}

	piiTags := "{}"
	if bq.PiiTags.RawMessage != nil {
		piiTags = string(bq.PiiTags.RawMessage)
	}

	return models.BigQuery{
		ID:            bq.ID,
		DatasetID:     bq.DatasetID,
		ProjectID:     bq.ProjectID,
		Dataset:       bq.Dataset,
		Table:         bq.TableName,
		TableType:     models.BigQueryType(strings.ToLower(bq.TableType)),
		LastModified:  bq.LastModified,
		Created:       bq.Created,
		Expires:       nullTimeToPtr(bq.Expires),
		Description:   bq.Description.String,
		PiiTags:       &piiTags,
		MissingSince:  &bq.MissingSince.Time,
		PseudoColumns: bq.PseudoColumns,
	}, nil
}

func (r *Repo) UpdateBigqueryDatasource(ctx context.Context, id uuid.UUID, schema json.RawMessage,
	lastModified, expires time.Time, description string, pseudoColumns []string,
) error {
	err := r.querier.UpdateBigqueryDatasourceSchema(ctx, gensql.UpdateBigqueryDatasourceSchemaParams{
		DatasetID: id,
		Schema: pqtype.NullRawMessage{
			RawMessage: schema,
			Valid:      true,
		},
		LastModified:  lastModified,
		Expires:       sql.NullTime{Time: expires, Valid: !expires.IsZero()},
		Description:   sql.NullString{String: description, Valid: true},
		PseudoColumns: pseudoColumns,
	})
	if err != nil {
		return fmt.Errorf("updating datasource_bigquery schema: %w", err)
	}

	return nil
}

func (r *Repo) UpdateBigqueryDatasourceMissing(ctx context.Context, datasetID uuid.UUID) error {
	return r.querier.UpdateBigqueryDatasourceMissing(ctx, datasetID)
}

func (r *Repo) GetDatasetMetadata(ctx context.Context, id uuid.UUID) ([]*models.TableColumn, error) {
	fmt.Println("database getdatasetmetadata")
	ds, err := r.querier.GetBigqueryDatasource(ctx, gensql.GetBigqueryDatasourceParams{
		DatasetID:   id,
		IsReference: false,
	})
	if err != nil {
		return nil, fmt.Errorf("getting bigquery datasource from database: %w", err)
	}

	var schema []*models.TableColumn
	if ds.Schema.Valid {
		if err := json.Unmarshal(ds.Schema.RawMessage, &schema); err != nil {
			return nil, fmt.Errorf("unmarshalling schema: %w", err)
		}
	}

	return schema, nil
}

func (r *Repo) GetDatasetPiiTags(ctx context.Context, id uuid.UUID) (map[string]string, error) {
	fmt.Println("database getdatasetpiitags")
	ds, err := r.querier.GetBigqueryDatasource(ctx, gensql.GetBigqueryDatasourceParams{
		DatasetID:   id,
		IsReference: false,
	})
	if err != nil {
		return nil, fmt.Errorf("getting bigquery datasource from database: %w", err)
	}

	piiTags := make(map[string]string)
	err = json.Unmarshal(ds.PiiTags.RawMessage, &piiTags)
	if err != nil {
		return nil, err
	}

	return piiTags, nil
}

func (r *Repo) GetDatasetsByUserAccess(ctx context.Context, user string) ([]*models.Dataset, error) {
	res, err := r.querier.GetDatasetsByUserAccess(ctx, user)
	if err != nil {
		return nil, err
	}

	return datasetsFromSQL(res)
}

func (r *Repo) GetDatasetsForOwner(ctx context.Context, userGroups []string) ([]*models.Dataset, error) {
	dprows, err := r.querier.GetDatasetsForOwner(ctx, userGroups)
	if err != nil {
		return nil, err
	}
	return datasetsFromSQL(dprows)
}

func (r *Repo) DeleteDataset(ctx context.Context, id uuid.UUID) error {
	r.events.TriggerDatasetDelete(ctx, id)

	if err := r.querier.DeleteDataset(ctx, id); err != nil {
		return fmt.Errorf("deleting dataset from database: %w", err)
	}

	return nil
}

func (r *Repo) GetAccessiblePseudoDatasourcesByUser(ctx context.Context, subjectsAsOwner []string, subjectsAsAccesser []string) ([]*models.PseudoDataset, error) {
	rows, err := r.querier.GetAccessiblePseudoDatasetsByUser(ctx, gensql.GetAccessiblePseudoDatasetsByUserParams{
		OwnerSubjects:  subjectsAsOwner,
		AccessSubjects: subjectsAsAccesser,
	})
	if err != nil {
		return nil, err
	}

	pseudoDatasets := []*models.PseudoDataset{}
	bqIDMap := make(map[string]int)
	for _, d := range rows {
		pseudoDataset, bqID := PseudoDatasetFromSQL(&d)
		_, exist := bqIDMap[bqID]
		if exist {
			continue
		}
		bqIDMap[bqID] = 1
		pseudoDatasets = append(pseudoDatasets, pseudoDataset)
	}
	return pseudoDatasets, nil
}

func (r *Repo) GetPseudoDatasourcesToDelete(ctx context.Context) ([]*models.BigQuery, error) {
	rows, err := r.querier.GetPseudoDatasourcesToDelete(ctx)
	if err != nil {
		return nil, err
	}

	pseudoViews := []*models.BigQuery{}
	for _, d := range rows {
		pseudoViews = append(pseudoViews, &models.BigQuery{
			ID:            d.ID,
			Dataset:       d.Dataset,
			ProjectID:     d.ProjectID,
			Table:         d.TableName,
			PseudoColumns: d.PseudoColumns,
		})
	}
	return pseudoViews, nil
}

func (r *Repo) SetDatasourceDeleted(ctx context.Context, id uuid.UUID) error {
	return r.querier.SetDatasourceDeleted(ctx, id)
}

func (r *Repo) GetOwnerGroupOfDataset(ctx context.Context, datasetID uuid.UUID) (string, error) {
	return r.querier.GetOwnerGroupOfDataset(ctx, datasetID)
}

func PseudoDatasetFromSQL(d *gensql.GetAccessiblePseudoDatasetsByUserRow) (*models.PseudoDataset, string) {
	return &models.PseudoDataset{
		// name is the name of the dataset
		Name: d.Name,
		// datasetID is the id of the dataset
		DatasetID: d.DatasetID,
		// datasourceID is the id of the bigquery datasource
		DatasourceID: d.BqDatasourceID,
	}, fmt.Sprintf("%v.%v.%v", d.BqProjectID, d.BqDatasetID, d.BqTableID)
}

func minimalDatasetFromSQL(ds gensql.Dataset) *models.Dataset {
	return &models.Dataset{
		ID:                       ds.ID,
		Name:                     ds.Name,
		Created:                  ds.Created,
		LastModified:             ds.LastModified,
		Description:              nullStringToPtr(ds.Description),
		Slug:                     ds.Slug,
		Repo:                     nullStringToPtr(ds.Repo),
		Pii:                      models.PiiLevel(ds.Pii),
		Keywords:                 ds.Keywords,
		Type:                     ds.Type,
		DataproductID:            ds.DataproductID,
		AnonymisationDescription: nullStringToPtr(ds.AnonymisationDescription),
		TargetUser:               nullStringToPtr(ds.TargetUser),
	}
}

func datasetsFromSQL(dsrows []gensql.DataproductCompleteView) ([]*models.Dataset, error) {
	datasets := []*models.Dataset{}

	for _, dsrow := range dsrows {
		owner := &models.Owner{
			Group:            dsrow.DpGroup,
			TeamkatalogenURL: nullStringToPtr(dsrow.TeamkatalogenUrl),
			TeamContact:      nullStringToPtr(dsrow.TeamContact),
			TeamID:           nullStringToPtr(dsrow.TeamID),
		}

		if !dsrow.DsID.Valid {
			continue
		}

		piiTags := "{}"
		if dsrow.PiiTags.RawMessage != nil {
			piiTags = string(dsrow.PiiTags.RawMessage)
		}

		var ds *models.Dataset

		for _, dsIn := range datasets {
			if dsIn.ID == dsrow.DsID.UUID {
				ds = dsIn
				break
			}
		}
		if ds == nil {
			ds = &models.Dataset{
				ID:            dsrow.DsID.UUID,
				Name:          dsrow.DsName.String,
				Created:       dsrow.DsCreated.Time,
				LastModified:  dsrow.DsLastModified.Time,
				Description:   nullStringToPtr(dsrow.DsDescription),
				Slug:          dsrow.DsSlug.String,
				Keywords:      dsrow.DsKeywords,
				DataproductID: dsrow.DataproductID,
				Owner:         owner,
				Mappings:      []models.MappingService{},
				Access:        []*models.Access{},
				Services:      &models.DatasetServices{},
			}
			datasets = append(datasets, ds)
		}

		if dsrow.BqID != uuid.Nil {
			var schema []*models.TableColumn
			if dsrow.BqSchema.Valid {
				if err := json.Unmarshal(dsrow.BqSchema.RawMessage, &schema); err != nil {
					return nil, fmt.Errorf("unmarshalling schema: %w", err)
				}
			}

			dsrc := models.BigQuery{
				ID:            dsrow.BqID,
				DatasetID:     dsrow.DsID.UUID,
				ProjectID:     dsrow.BqProject,
				Dataset:       dsrow.BqDataset,
				Table:         dsrow.BqTableName,
				TableType:     models.BigQueryType(dsrow.BqTableType),
				Created:       dsrow.BqCreated,
				LastModified:  dsrow.BqLastModified,
				Expires:       nullTimeToPtr(dsrow.BqExpires),
				Description:   dsrow.BqDescription.String,
				PiiTags:       &piiTags,
				MissingSince:  nullTimeToPtr(dsrow.BqMissingSince),
				PseudoColumns: dsrow.PseudoColumns,
				Schema:        schema,
			}
			ds.Datasource = dsrc
		}

		if len(dsrow.MappingServices) > 0 {
			for _, service := range dsrow.MappingServices {
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

		if dsrow.AccessID.Valid {
			exist := false
			for _, dsAccess := range ds.Access {
				if dsAccess.ID == dsrow.AccessID.UUID {
					exist = true
					break
				}
			}
			if !exist {
				access := &models.Access{
					ID:              dsrow.AccessID.UUID,
					Subject:         dsrow.AccessSubject.String,
					Granter:         dsrow.AccessGranter.String,
					Expires:         nullTimeToPtr(dsrow.AccessExpires),
					Created:         dsrow.AccessCreated.Time,
					Revoked:         nullTimeToPtr(dsrow.AccessRevoked),
					DatasetID:       dsrow.DsID.UUID,
					AccessRequestID: nullUUIDToUUIDPtr(dsrow.AccessRequestID),
				}
				ds.Access = append(ds.Access, access)
			}
		}

		if ds.Services == nil && dsrow.MbDatabaseID.Valid {
			svc := &models.DatasetServices{}
			base := "https://metabase.intern.dev.nav.no/browse/%v"
			if os.Getenv("NAIS_CLUSTER_NAME") == "prod-gcp" {
				base = "https://metabase.intern.nav.no/browse/%v"
			}
			url := fmt.Sprintf(base, dsrow.MbDatabaseID.Int32)
			svc.Metabase = &url
			ds.Services = svc
		}
	}

	return datasets, nil
}
