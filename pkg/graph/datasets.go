package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/navikt/nada-backend/pkg/graph/models"
)

func getPIIColumns(piiTags *string) ([]string, error) {
	piiTagsBytes := []byte(*piiTags)

	tagMap := map[string]string{}

	if err := json.Unmarshal(piiTagsBytes, &tagMap); err != nil {
		return nil, err
	}

	piiColumns := []string{}
	for key, value := range tagMap {
		//TODO: fix the magic string
		if value == "PII_DirekteIdentifiserende" {
			piiColumns = append(piiColumns, key)
		}
	}

	return piiColumns, nil
}

func (r *mutationResolver) CreatePseudoynimizedView(ctx context.Context, input models.NewDataset) (*models.NewDataset, error) {
	if !input.CreatePseudoynimizedView {
		return nil, nil
	}

	piiColumns, err := getPIIColumns(input.BigQuery.PiiTags)
	if err != nil {
		r.log.WithError(err).Errorf("failed to parse PII columns for input dataset %v: %v", input.Name, input.BigQuery.PiiTags)
		return nil, err
	}

	if len(piiColumns) == 0 {
		return nil, nil
	}

	project, dataset, view, err := r.bigquery.CreatePseudoynimizedView(ctx, input.BigQuery.ProjectID, input.BigQuery.Dataset, input.BigQuery.Table, piiColumns)
	if err != nil {
		r.log.WithError(err).Errorf("failed to create pseudoynimized view for dataset %v", input.Name)
		return nil, err
	}

	metadata, err := r.bigquery.TableMetadata(ctx, project, dataset, view)
	if err != nil {
		return &models.NewDataset{}, fmt.Errorf("failed to fetch metadata on pseudoynimized view %v in %v.%v",
			view, project, dataset)
	}

	if err := r.accessMgr.MakeAuthorizedViewForDataset(ctx, input.BigQuery.ProjectID, input.BigQuery.Dataset,
		project, dataset, view); err != nil {
		return &models.NewDataset{}, err
	}

	description := "Pseudonymisert versjon av " + input.Name + ". \nGenerert av markedplassen"
	return &models.NewDataset{
		DataproductID: input.DataproductID,
		Name:          input.Name + " (pseudonymisert)",
		Description:   &description,
		Pii:           input.Pii,
		Keywords:      input.Keywords,
		BigQuery: models.NewBigQuery{
			ProjectID: project,
			Dataset:   dataset,
			Table:     view,
			PiiTags:   input.BigQuery.PiiTags,
		},
		AnonymisationDescription: input.AnonymisationDescription,
		GrantAllUsers:            input.GrantAllUsers,
		TargetUser:               input.TargetUser,
		CreatePseudoynimizedView: input.CreatePseudoynimizedView,
		Metadata:                 metadata,
	}, nil
}
