package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
)

// TriggerMetadataSync is the resolver for the triggerMetadataSync field.
func (r *mutationResolver) TriggerMetadataSync(ctx context.Context) (bool, error) {
	bqs, err := r.repo.GetBigqueryDatasources(ctx)
	if err != nil {
		return false, err
	}

	var errs errorList

	for _, bq := range bqs {
		err := r.UpdateMetadata(ctx, bq)
		if err != nil {
			errs = r.handleSyncError(ctx, errs, err, bq)
		}
	}
	if len(errs) != 0 {
		return false, errs
	}
	return true, nil
}
