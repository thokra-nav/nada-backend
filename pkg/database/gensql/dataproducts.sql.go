// Code generated by sqlc. DO NOT EDIT.
// source: dataproducts.sql

package gensql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tabbed/pqtype"
)

const createBigqueryDatasource = `-- name: CreateBigqueryDatasource :one
INSERT INTO datasource_bigquery ("dataproduct_id",
                                 "project_id",
                                 "dataset",
                                 "table_name",
                                 "schema",
                                 "last_modified",
                                 "created",
                                 "expires",
                                 "table_type")
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9)
RETURNING dataproduct_id, project_id, dataset, table_name, schema, last_modified, created, expires, table_type
`

type CreateBigqueryDatasourceParams struct {
	DataproductID uuid.UUID
	ProjectID     string
	Dataset       string
	TableName     string
	Schema        pqtype.NullRawMessage
	LastModified  time.Time
	Created       time.Time
	Expires       sql.NullTime
	TableType     string
}

func (q *Queries) CreateBigqueryDatasource(ctx context.Context, arg CreateBigqueryDatasourceParams) (DatasourceBigquery, error) {
	row := q.db.QueryRowContext(ctx, createBigqueryDatasource,
		arg.DataproductID,
		arg.ProjectID,
		arg.Dataset,
		arg.TableName,
		arg.Schema,
		arg.LastModified,
		arg.Created,
		arg.Expires,
		arg.TableType,
	)
	var i DatasourceBigquery
	err := row.Scan(
		&i.DataproductID,
		&i.ProjectID,
		&i.Dataset,
		&i.TableName,
		&i.Schema,
		&i.LastModified,
		&i.Created,
		&i.Expires,
		&i.TableType,
	)
	return i, err
}

const createDataproduct = `-- name: CreateDataproduct :one
INSERT INTO dataproducts ("name",
                          "description",
                          "pii",
                          "type",
                          "group",
                          "teamkatalogen_url",
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
        $8,
        $9)
RETURNING id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
`

type CreateDataproductParams struct {
	Name                  string
	Description           sql.NullString
	Pii                   bool
	Type                  DatasourceType
	OwnerGroup            string
	OwnerTeamkatalogenUrl sql.NullString
	Slug                  string
	Repo                  sql.NullString
	Keywords              []string
}

func (q *Queries) CreateDataproduct(ctx context.Context, arg CreateDataproductParams) (Dataproduct, error) {
	row := q.db.QueryRowContext(ctx, createDataproduct,
		arg.Name,
		arg.Description,
		arg.Pii,
		arg.Type,
		arg.OwnerGroup,
		arg.OwnerTeamkatalogenUrl,
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
		&i.TeamkatalogenUrl,
	)
	return i, err
}

const createDataproductRequester = `-- name: CreateDataproductRequester :exec
INSERT INTO dataproduct_requesters (dataproduct_id, "subject")
VALUES ($1, LOWER($2))
`

type CreateDataproductRequesterParams struct {
	DataproductID uuid.UUID
	Subject       string
}

func (q *Queries) CreateDataproductRequester(ctx context.Context, arg CreateDataproductRequesterParams) error {
	_, err := q.db.ExecContext(ctx, createDataproductRequester, arg.DataproductID, arg.Subject)
	return err
}

const dataproductGroupStats = `-- name: DataproductGroupStats :many
SELECT "group",
       count(1) as "count"
FROM "dataproducts"
GROUP BY "group"
ORDER BY "count" DESC
LIMIT $2 OFFSET $1
`

type DataproductGroupStatsParams struct {
	Offs int32
	Lim  int32
}

type DataproductGroupStatsRow struct {
	Group string
	Count int64
}

func (q *Queries) DataproductGroupStats(ctx context.Context, arg DataproductGroupStatsParams) ([]DataproductGroupStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, dataproductGroupStats, arg.Offs, arg.Lim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DataproductGroupStatsRow{}
	for rows.Next() {
		var i DataproductGroupStatsRow
		if err := rows.Scan(&i.Group, &i.Count); err != nil {
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

const dataproductKeywords = `-- name: DataproductKeywords :many
SELECT keyword::text, count(1) as "count"
FROM (
	SELECT unnest(keywords) as keyword
	FROM dataproducts
) s
WHERE true
AND CASE WHEN coalesce(TRIM($1), '') = '' THEN true ELSE keyword ILIKE $1::text || '%' END
GROUP BY keyword
ORDER BY "count" DESC
LIMIT 15
`

type DataproductKeywordsRow struct {
	Keyword string
	Count   int64
}

func (q *Queries) DataproductKeywords(ctx context.Context, keyword string) ([]DataproductKeywordsRow, error) {
	rows, err := q.db.QueryContext(ctx, dataproductKeywords, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DataproductKeywordsRow{}
	for rows.Next() {
		var i DataproductKeywordsRow
		if err := rows.Scan(&i.Keyword, &i.Count); err != nil {
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

const dataproductsByMetabase = `-- name: DataproductsByMetabase :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
FROM dataproducts
WHERE id IN (
	SELECT dataproduct_id
	FROM metabase_metadata
)
ORDER BY last_modified DESC
LIMIT $2 OFFSET $1
`

type DataproductsByMetabaseParams struct {
	Offs int32
	Lim  int32
}

func (q *Queries) DataproductsByMetabase(ctx context.Context, arg DataproductsByMetabaseParams) ([]Dataproduct, error) {
	rows, err := q.db.QueryContext(ctx, dataproductsByMetabase, arg.Offs, arg.Lim)
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
			&i.TeamkatalogenUrl,
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

const deleteDataproduct = `-- name: DeleteDataproduct :exec
DELETE
FROM dataproducts
WHERE id = $1
`

func (q *Queries) DeleteDataproduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteDataproduct, id)
	return err
}

const deleteDataproductRequester = `-- name: DeleteDataproductRequester :exec
DELETE
FROM dataproduct_requesters
WHERE dataproduct_id = $1
  AND "subject" = LOWER($2)
`

type DeleteDataproductRequesterParams struct {
	DataproductID uuid.UUID
	Subject       string
}

func (q *Queries) DeleteDataproductRequester(ctx context.Context, arg DeleteDataproductRequesterParams) error {
	_, err := q.db.ExecContext(ctx, deleteDataproductRequester, arg.DataproductID, arg.Subject)
	return err
}

const getBigqueryDatasource = `-- name: GetBigqueryDatasource :one
SELECT dataproduct_id, project_id, dataset, table_name, schema, last_modified, created, expires, table_type
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
		&i.LastModified,
		&i.Created,
		&i.Expires,
		&i.TableType,
	)
	return i, err
}

const getBigqueryDatasources = `-- name: GetBigqueryDatasources :many
SELECT dataproduct_id, project_id, dataset, table_name, schema, last_modified, created, expires, table_type
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
			&i.LastModified,
			&i.Created,
			&i.Expires,
			&i.TableType,
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
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
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
		&i.TeamkatalogenUrl,
	)
	return i, err
}

const getDataproductRequesters = `-- name: GetDataproductRequesters :many
SELECT "subject"
FROM dataproduct_requesters
WHERE dataproduct_id = $1
`

func (q *Queries) GetDataproductRequesters(ctx context.Context, dataproductID uuid.UUID) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getDataproductRequesters, dataproductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var subject string
		if err := rows.Scan(&subject); err != nil {
			return nil, err
		}
		items = append(items, subject)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDataproducts = `-- name: GetDataproducts :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
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
			&i.TeamkatalogenUrl,
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

const getDataproductsByGroups = `-- name: GetDataproductsByGroups :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
FROM dataproducts
WHERE "group" = ANY ($1::text[])
ORDER BY last_modified DESC
`

func (q *Queries) GetDataproductsByGroups(ctx context.Context, groups []string) ([]Dataproduct, error) {
	rows, err := q.db.QueryContext(ctx, getDataproductsByGroups, pq.Array(groups))
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
			&i.TeamkatalogenUrl,
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

const getDataproductsByIDs = `-- name: GetDataproductsByIDs :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
FROM dataproducts
WHERE id = ANY ($1::uuid[])
ORDER BY last_modified DESC
`

func (q *Queries) GetDataproductsByIDs(ctx context.Context, ids []uuid.UUID) ([]Dataproduct, error) {
	rows, err := q.db.QueryContext(ctx, getDataproductsByIDs, pq.Array(ids))
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
			&i.TeamkatalogenUrl,
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

const getDataproductsByUserAccess = `-- name: GetDataproductsByUserAccess :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
FROM dataproducts
WHERE id = ANY (SELECT dataproduct_id
                FROM dataproduct_access
                WHERE "subject" = LOWER($1)
                  AND revoked IS NULL
                  AND (expires > NOW() OR expires IS NULL))
ORDER BY last_modified DESC
`

func (q *Queries) GetDataproductsByUserAccess(ctx context.Context, id string) ([]Dataproduct, error) {
	rows, err := q.db.QueryContext(ctx, getDataproductsByUserAccess, id)
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
			&i.TeamkatalogenUrl,
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
SET "schema"        = $1,
    "last_modified" = $2,
    "expires"       = $3
WHERE dataproduct_id = $4
`

type UpdateBigqueryDatasourceSchemaParams struct {
	Schema        pqtype.NullRawMessage
	LastModified  time.Time
	Expires       sql.NullTime
	DataproductID uuid.UUID
}

func (q *Queries) UpdateBigqueryDatasourceSchema(ctx context.Context, arg UpdateBigqueryDatasourceSchemaParams) error {
	_, err := q.db.ExecContext(ctx, updateBigqueryDatasourceSchema,
		arg.Schema,
		arg.LastModified,
		arg.Expires,
		arg.DataproductID,
	)
	return err
}

const updateDataproduct = `-- name: UpdateDataproduct :one
UPDATE dataproducts
SET "name"              = $1,
    "description"       = $2,
    "pii"               = $3,
    "slug"              = $4,
    "repo"              = $5,
    "teamkatalogen_url" = $6,
    "keywords"          = $7
WHERE id = $8
RETURNING id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
`

type UpdateDataproductParams struct {
	Name                  string
	Description           sql.NullString
	Pii                   bool
	Slug                  string
	Repo                  sql.NullString
	OwnerTeamkatalogenUrl sql.NullString
	Keywords              []string
	ID                    uuid.UUID
}

func (q *Queries) UpdateDataproduct(ctx context.Context, arg UpdateDataproductParams) (Dataproduct, error) {
	row := q.db.QueryRowContext(ctx, updateDataproduct,
		arg.Name,
		arg.Description,
		arg.Pii,
		arg.Slug,
		arg.Repo,
		arg.OwnerTeamkatalogenUrl,
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
		&i.TeamkatalogenUrl,
	)
	return i, err
}
