package access

import "context"

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (a Noop) Grant(ctx context.Context, projectID, datasetID, tableID, member string) error {
	return nil
}

func (a Noop) Revoke(ctx context.Context, projectID, datasetID, tableID, member string) error {
	return nil
}

func (a Noop) HasAccess(ctx context.Context, projectID, datasetID, tableID, member string) (bool, error) {
	return true, nil
}

func (a Noop) AddToAuthorizedViews(ctx context.Context, projectID, dataset, table string) error {
	return nil
}

func (a Noop) MakeAuthorizedViewForDataset(ctx context.Context, projectID, dataset, viewProjectID, viewDataset, viewID string) error {
	return nil
}
