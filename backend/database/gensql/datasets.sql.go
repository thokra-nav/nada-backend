// Code generated by sqlc. DO NOT EDIT.
// source: datasets.sql

package gensql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createDataset = `-- name: CreateDataset :one
INSERT INTO datasets (
	"dataproduct_id",
	"name",
	"description",
	"pii",
	"project_id",
	"dataset",
	"table_name",
	"type"
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
) RETURNING id, dataproduct_id, name, description, pii, created, last_modified, project_id, dataset, table_name, type
`

type CreateDatasetParams struct {
	DataproductID uuid.UUID
	Name          string
	Description   sql.NullString
	Pii           bool
	ProjectID     string
	Dataset       string
	TableName     string
	Type          string
}

func (q *Queries) CreateDataset(ctx context.Context, arg CreateDatasetParams) (Dataset, error) {
	row := q.db.QueryRowContext(ctx, createDataset,
		arg.DataproductID,
		arg.Name,
		arg.Description,
		arg.Pii,
		arg.ProjectID,
		arg.Dataset,
		arg.TableName,
		arg.Type,
	)
	var i Dataset
	err := row.Scan(
		&i.ID,
		&i.DataproductID,
		&i.Name,
		&i.Description,
		&i.Pii,
		&i.Created,
		&i.LastModified,
		&i.ProjectID,
		&i.Dataset,
		&i.TableName,
		&i.Type,
	)
	return i, err
}

const deleteDataset = `-- name: DeleteDataset :exec
DELETE FROM datasets WHERE id = $1
`

func (q *Queries) DeleteDataset(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteDataset, id)
	return err
}

const getDatasets = `-- name: GetDatasets :many
SELECT id, dataproduct_id, name, description, pii, created, last_modified, project_id, dataset, table_name, type FROM datasets WHERE dataproduct_id = $1
`

func (q *Queries) GetDatasets(ctx context.Context, dataproductID uuid.UUID) ([]Dataset, error) {
	rows, err := q.db.QueryContext(ctx, getDatasets, dataproductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Dataset{}
	for rows.Next() {
		var i Dataset
		if err := rows.Scan(
			&i.ID,
			&i.DataproductID,
			&i.Name,
			&i.Description,
			&i.Pii,
			&i.Created,
			&i.LastModified,
			&i.ProjectID,
			&i.Dataset,
			&i.TableName,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateDataset = `-- name: UpdateDataset :one
UPDATE datasets SET
	"dataproduct_id" = $1,
	"name" = $2,
	"description" = $3,
	"pii" = $4
WHERE id = $5
RETURNING id, dataproduct_id, name, description, pii, created, last_modified, project_id, dataset, table_name, type
`

type UpdateDatasetParams struct {
	DataproductID uuid.UUID
	Name          string
	Description   sql.NullString
	Pii           bool
	ID            uuid.UUID
}

func (q *Queries) UpdateDataset(ctx context.Context, arg UpdateDatasetParams) (Dataset, error) {
	row := q.db.QueryRowContext(ctx, updateDataset,
		arg.DataproductID,
		arg.Name,
		arg.Description,
		arg.Pii,
		arg.ID,
	)
	var i Dataset
	err := row.Scan(
		&i.ID,
		&i.DataproductID,
		&i.Name,
		&i.Description,
		&i.Pii,
		&i.Created,
		&i.LastModified,
		&i.ProjectID,
		&i.Dataset,
		&i.TableName,
		&i.Type,
	)
	return i, err
}
