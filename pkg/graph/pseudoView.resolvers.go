package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"context"
	"fmt"

	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

// CreatePseudoView is the resolver for the createPseudoView field.
func (r *mutationResolver) CreatePseudoView(ctx context.Context, input models.NewPseudoView) (string, error) {
	pid, did, tid, err := r.bigquery.CreatePseudonymisedView(ctx, input.ProjectID, input.Dataset, input.Table, input.TargetColumns)
	if err != nil {
		r.log.WithError(err).Errorf("failed to create pseudonymised view for %v %v %v", input.ProjectID, input.Dataset, input.Table)
	}
	return fmt.Sprintf("%v.%v.%v", pid, did, tid), err
}

// CreateJoinableViews is the resolver for the createJoinableViews field.
func (r *mutationResolver) CreateJoinableViews(ctx context.Context, input models.NewJoinableViews) (string, error) {
	user := auth.GetUser(ctx)
	datasets := []*models.Dataset{}
	for _, dsid := range input.DatasetIDs {
		var dataset *models.Dataset
		dataset, err := r.repo.GetDataset(ctx, dsid)
		if err != nil {
			return "", fmt.Errorf("Failed to find dataset to make joinable view: %v", err)
		}
		dataproduct, err := r.repo.GetDataproduct(ctx, dataset.DataproductID)
		if err != nil {
			return "", fmt.Errorf("Failed to find dataproduct for dataset: %v", err)
		}
		if !user.GoogleGroups.Contains(dataproduct.Owner.Group) {
			access, err := r.repo.ListActiveAccessToDataset(ctx, dataset.ID)
			if err != nil {
				return "", fmt.Errorf("Failed to check dataset access: %v", err)
			}
			accessSet := make(map[string]int)
			for _, da := range access {
				accessSet[da.Subject] = 1
			}
			for _, ugg := range user.GoogleGroups {
				accessSet["group:"+ugg.Email] = 1
			}
			accessSet["user:"+user.Email] = 1
			if len(accessSet) == len(user.GoogleGroups.Emails())+1+len(access) {
				return "", fmt.Errorf("Access denied")
			}
		}
		datasets = append(datasets, dataset)
	}

	tableUrls := []models.BigQuery{}
	for _, ds := range datasets {
		if datasource, err := r.repo.GetBigqueryDatasource(ctx, ds.ID); err == nil {
			tableUrls = append(tableUrls, datasource)
		} else {
			return "", fmt.Errorf("Failed to find bigquery datasource: %v", err)
		}
	}

	projectID, joinableDatasetID, views, err := r.bigquery.CreateJoinableViewsForUser(ctx, auth.GetUser(ctx), tableUrls)
	if err != nil {
		return "", err
	}

	for _, dstbl := range tableUrls {
		if err := r.accessMgr.AddToAuthorizedViews(ctx, dstbl.ProjectID, dstbl.Dataset, projectID, joinableDatasetID, dstbl.Table); err != nil {
			return "", fmt.Errorf("Failed to add to authorized views: %v", err)
		}
		if err := r.accessMgr.AddToAuthorizedViews(ctx, projectID, "secrets_vault", projectID, joinableDatasetID, dstbl.Table); err != nil {
			return "", fmt.Errorf("Failed to add to secrets' authorized views: %v", err)
		}
	}

	subj := user.Email
	subjType := models.SubjectTypeUser
	subjWithType := subjType.String() + ":" + subj

	for _, v := range views {
		if err := r.accessMgr.Grant(ctx, projectID, joinableDatasetID, v, subjWithType); err != nil {
			return "", err
		}

	}

	return joinableDatasetID, nil
}

// JoinableViews is the resolver for the joinableViews field.
func (r *queryResolver) JoinableViews(ctx context.Context) ([]*models.JoinableView, error) {
	return r.bigquery.GetJoinableViewsForUser(ctx, auth.GetUser(ctx))
}
