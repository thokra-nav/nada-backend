package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

func (r *Repo) GetPolly(ctx context.Context, accessID uuid.UUID) (*models.Polly, error) {
	pollySQL, err := r.querier.GetPolly(ctx, accessID)
	if err != nil {
		return nil, err
	}

	return &models.Polly{
		ID:   pollySQL.PollyID,
		Name: pollySQL.PollyName,
		URL:  pollySQL.PollyUrl,
	}, nil
}
