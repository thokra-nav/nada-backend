// Code generated by sqlc. DO NOT EDIT.
// source: metabase_metadata.sql

package gensql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createMetabaseMetadata = `-- name: CreateMetabaseMetadata :exec
INSERT INTO metabase_metadata (
    "dataset_id",
    "database_id",
    "permission_group_id",
    "collection_id",
    "sa_email",
    "deleted_at"
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
`

type CreateMetabaseMetadataParams struct {
	DatasetID         uuid.UUID
	DatabaseID        int32
	PermissionGroupID sql.NullInt32
	CollectionID      sql.NullInt32
	SaEmail           string
	DeletedAt         sql.NullTime
}

func (q *Queries) CreateMetabaseMetadata(ctx context.Context, arg CreateMetabaseMetadataParams) error {
	_, err := q.db.ExecContext(ctx, createMetabaseMetadata,
		arg.DatasetID,
		arg.DatabaseID,
		arg.PermissionGroupID,
		arg.CollectionID,
		arg.SaEmail,
		arg.DeletedAt,
	)
	return err
}

const deleteMetabaseMetadata = `-- name: DeleteMetabaseMetadata :exec
DELETE 
FROM metabase_metadata
WHERE "dataset_id" = $1
`

func (q *Queries) DeleteMetabaseMetadata(ctx context.Context, datasetID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteMetabaseMetadata, datasetID)
	return err
}

const getMetabaseMetadata = `-- name: GetMetabaseMetadata :one
SELECT database_id, permission_group_id, sa_email, collection_id, deleted_at, dataset_id
FROM metabase_metadata
WHERE "dataset_id" = $1 AND "deleted_at" IS NULL
`

func (q *Queries) GetMetabaseMetadata(ctx context.Context, datasetID uuid.UUID) (MetabaseMetadatum, error) {
	row := q.db.QueryRowContext(ctx, getMetabaseMetadata, datasetID)
	var i MetabaseMetadatum
	err := row.Scan(
		&i.DatabaseID,
		&i.PermissionGroupID,
		&i.SaEmail,
		&i.CollectionID,
		&i.DeletedAt,
		&i.DatasetID,
	)
	return i, err
}

const getMetabaseMetadataWithDeleted = `-- name: GetMetabaseMetadataWithDeleted :one
SELECT database_id, permission_group_id, sa_email, collection_id, deleted_at, dataset_id
FROM metabase_metadata
WHERE "dataset_id" = $1
`

func (q *Queries) GetMetabaseMetadataWithDeleted(ctx context.Context, datasetID uuid.UUID) (MetabaseMetadatum, error) {
	row := q.db.QueryRowContext(ctx, getMetabaseMetadataWithDeleted, datasetID)
	var i MetabaseMetadatum
	err := row.Scan(
		&i.DatabaseID,
		&i.PermissionGroupID,
		&i.SaEmail,
		&i.CollectionID,
		&i.DeletedAt,
		&i.DatasetID,
	)
	return i, err
}

const restoreMetabaseMetadata = `-- name: RestoreMetabaseMetadata :exec
UPDATE metabase_metadata
SET "deleted_at" = null
WHERE dataset_id = $1
`

func (q *Queries) RestoreMetabaseMetadata(ctx context.Context, datasetID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, restoreMetabaseMetadata, datasetID)
	return err
}

const setPermissionGroupMetabaseMetadata = `-- name: SetPermissionGroupMetabaseMetadata :exec
UPDATE metabase_metadata
SET "permission_group_id" = $1
WHERE dataset_id = $2
`

type SetPermissionGroupMetabaseMetadataParams struct {
	ID        sql.NullInt32
	DatasetID uuid.UUID
}

func (q *Queries) SetPermissionGroupMetabaseMetadata(ctx context.Context, arg SetPermissionGroupMetabaseMetadataParams) error {
	_, err := q.db.ExecContext(ctx, setPermissionGroupMetabaseMetadata, arg.ID, arg.DatasetID)
	return err
}

const softDeleteMetabaseMetadata = `-- name: SoftDeleteMetabaseMetadata :exec
UPDATE metabase_metadata
SET "deleted_at" = NOW()
WHERE dataset_id = $1
`

func (q *Queries) SoftDeleteMetabaseMetadata(ctx context.Context, datasetID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, softDeleteMetabaseMetadata, datasetID)
	return err
}
