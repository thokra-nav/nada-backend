// Code generated by sqlc. DO NOT EDIT.
// source: dataproducts.sql

package gensql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tabbed/pqtype"
)

const createBigqueryDatasource = `-- name: CreateBigqueryDatasource :one
INSERT INTO datasource_bigquery ("dataproduct_id",
                                 "project_id",
                                 "dataset",
                                 "table_name")
VALUES ($1,
        $2,
        $3,
        $4)
RETURNING dataproduct_id, project_id, dataset, table_name, schema
`

type CreateBigqueryDatasourceParams struct {
	DataproductID uuid.UUID
	ProjectID     string
	Dataset       string
	TableName     string
}

func (q *Queries) CreateBigqueryDatasource(ctx context.Context, arg CreateBigqueryDatasourceParams) (DatasourceBigquery, error) {
	row := q.db.QueryRowContext(ctx, createBigqueryDatasource,
		arg.DataproductID,
		arg.ProjectID,
		arg.Dataset,
		arg.TableName,
	)
	var i DatasourceBigquery
	err := row.Scan(
		&i.DataproductID,
		&i.ProjectID,
		&i.Dataset,
		&i.TableName,
		&i.Schema,
	)
	return i, err
}

const createDataproduct = `-- name: CreateDataproduct :one
INSERT INTO dataproducts ("name",
                          "description",
                          "pii",
                          "type",
                          "group",
                          "slug",
                          "repo",
                          "keywords")
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8)
RETURNING id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords
`

type CreateDataproductParams struct {
	Name        string
	Description sql.NullString
	Pii         bool
	Type        DatasourceType
	OwnerGroup  string
	Slug        sql.NullString
	Repo        sql.NullString
	Keywords    []string
}

func (q *Queries) CreateDataproduct(ctx context.Context, arg CreateDataproductParams) (Dataproduct, error) {
	row := q.db.QueryRowContext(ctx, createDataproduct,
		arg.Name,
		arg.Description,
		arg.Pii,
		arg.Type,
		arg.OwnerGroup,
		arg.Slug,
		arg.Repo,
		pq.Array(arg.Keywords),
	)
	var i Dataproduct
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Group,
		&i.Pii,
		&i.Created,
		&i.LastModified,
		&i.Type,
		&i.TsvDocument,
		&i.Slug,
		&i.Repo,
		pq.Array(&i.Keywords),
	)
	return i, err
}

const deleteDataproduct = `-- name: DeleteDataproduct :exec
DELETE
FROM dataproducts
WHERE id = $1
`

func (q *Queries) DeleteDataproduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteDataproduct, id)
	return err
}

const getBigqueryDatasource = `-- name: GetBigqueryDatasource :one
SELECT dataproduct_id, project_id, dataset, table_name, schema
FROM datasource_bigquery
WHERE dataproduct_id = $1
`

func (q *Queries) GetBigqueryDatasource(ctx context.Context, dataproductID uuid.UUID) (DatasourceBigquery, error) {
	row := q.db.QueryRowContext(ctx, getBigqueryDatasource, dataproductID)
	var i DatasourceBigquery
	err := row.Scan(
		&i.DataproductID,
		&i.ProjectID,
		&i.Dataset,
		&i.TableName,
		&i.Schema,
	)
	return i, err
}

const getBigqueryDatasources = `-- name: GetBigqueryDatasources :many
SELECT dataproduct_id, project_id, dataset, table_name, schema
FROM datasource_bigquery
`

func (q *Queries) GetBigqueryDatasources(ctx context.Context) ([]DatasourceBigquery, error) {
	rows, err := q.db.QueryContext(ctx, getBigqueryDatasources)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DatasourceBigquery{}
	for rows.Next() {
		var i DatasourceBigquery
		if err := rows.Scan(
			&i.DataproductID,
			&i.ProjectID,
			&i.Dataset,
			&i.TableName,
			&i.Schema,
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

const getDataproduct = `-- name: GetDataproduct :one
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords
FROM dataproducts
WHERE id = $1
`

func (q *Queries) GetDataproduct(ctx context.Context, id uuid.UUID) (Dataproduct, error) {
	row := q.db.QueryRowContext(ctx, getDataproduct, id)
	var i Dataproduct
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Group,
		&i.Pii,
		&i.Created,
		&i.LastModified,
		&i.Type,
		&i.TsvDocument,
		&i.Slug,
		&i.Repo,
		pq.Array(&i.Keywords),
	)
	return i, err
}

const getDataproducts = `-- name: GetDataproducts :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords
FROM dataproducts
ORDER BY last_modified DESC
LIMIT $2 OFFSET $1
`

type GetDataproductsParams struct {
	Offset int32
	Limit  int32
}

func (q *Queries) GetDataproducts(ctx context.Context, arg GetDataproductsParams) ([]Dataproduct, error) {
	rows, err := q.db.QueryContext(ctx, getDataproducts, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Dataproduct{}
	for rows.Next() {
		var i Dataproduct
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Group,
			&i.Pii,
			&i.Created,
			&i.LastModified,
			&i.Type,
			&i.TsvDocument,
			&i.Slug,
			&i.Repo,
			pq.Array(&i.Keywords),
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

const updateBigqueryDatasourceSchema = `-- name: UpdateBigqueryDatasourceSchema :exec
UPDATE datasource_bigquery
SET "schema" = $1
WHERE dataproduct_id = $2
`

type UpdateBigqueryDatasourceSchemaParams struct {
	Schema        pqtype.NullRawMessage
	DataproductID uuid.UUID
}

func (q *Queries) UpdateBigqueryDatasourceSchema(ctx context.Context, arg UpdateBigqueryDatasourceSchemaParams) error {
	_, err := q.db.ExecContext(ctx, updateBigqueryDatasourceSchema, arg.Schema, arg.DataproductID)
	return err
}

const updateDataproduct = `-- name: UpdateDataproduct :one
UPDATE dataproducts
SET "name"        = $1,
    "description" = $2,
    "pii"         = $3,
    "slug"        = $4,
    "repo"        = $5,
    "keywords"    = $6 
WHERE id = $7
RETURNING id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords
`

type UpdateDataproductParams struct {
	Name        string
	Description sql.NullString
	Pii         bool
	Slug        sql.NullString
	Repo        sql.NullString
	Keywords    []string
	ID          uuid.UUID
}

func (q *Queries) UpdateDataproduct(ctx context.Context, arg UpdateDataproductParams) (Dataproduct, error) {
	row := q.db.QueryRowContext(ctx, updateDataproduct,
		arg.Name,
		arg.Description,
		arg.Pii,
		arg.Slug,
		arg.Repo,
		pq.Array(arg.Keywords),
		arg.ID,
	)
	var i Dataproduct
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Group,
		&i.Pii,
		&i.Created,
		&i.LastModified,
		&i.Type,
		&i.TsvDocument,
		&i.Slug,
		&i.Repo,
		pq.Array(&i.Keywords),
	)
	return i, err
}
