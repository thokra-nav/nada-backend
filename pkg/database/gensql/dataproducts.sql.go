// Code generated by sqlc. DO NOT EDIT.
// source: dataproducts.sql

package gensql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createDataproduct = `-- name: CreateDataproduct :one
INSERT INTO dataproducts ("name",
                          "description",
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
        $7)
RETURNING id, name, description, "group", created, last_modified, tsv_document, slug, repo, keywords, teamkatalogen_url
`

type CreateDataproductParams struct {
	Name                  string
	Description           sql.NullString
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
		&i.Created,
		&i.LastModified,
		&i.TsvDocument,
		&i.Slug,
		&i.Repo,
		pq.Array(&i.Keywords),
		&i.TeamkatalogenUrl,
	)
	return i, err
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

const deleteDataproduct = `-- name: DeleteDataproduct :exec
DELETE
FROM dataproducts
WHERE id = $1
`

func (q *Queries) DeleteDataproduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteDataproduct, id)
	return err
}

const getDataproduct = `-- name: GetDataproduct :one
SELECT id, name, description, "group", created, last_modified, tsv_document, slug, repo, keywords, teamkatalogen_url
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
		&i.Created,
		&i.LastModified,
		&i.TsvDocument,
		&i.Slug,
		&i.Repo,
		pq.Array(&i.Keywords),
		&i.TeamkatalogenUrl,
	)
	return i, err
}

const getDataproducts = `-- name: GetDataproducts :many
SELECT id, name, description, "group", created, last_modified, tsv_document, slug, repo, keywords, teamkatalogen_url
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
			&i.Created,
			&i.LastModified,
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
SELECT id, name, description, "group", created, last_modified, tsv_document, slug, repo, keywords, teamkatalogen_url
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
			&i.Created,
			&i.LastModified,
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
SELECT id, name, description, "group", created, last_modified, tsv_document, slug, repo, keywords, teamkatalogen_url
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
			&i.Created,
			&i.LastModified,
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

const updateDataproduct = `-- name: UpdateDataproduct :one
UPDATE dataproducts
SET "name"              = $1,
    "description"       = $2,
    "slug"              = $3,
    "repo"              = $4,
    "teamkatalogen_url" = $5,
    "keywords"          = $6
WHERE id = $7
RETURNING id, name, description, "group", created, last_modified, tsv_document, slug, repo, keywords, teamkatalogen_url
`

type UpdateDataproductParams struct {
	Name                  string
	Description           sql.NullString
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
		&i.Created,
		&i.LastModified,
		&i.TsvDocument,
		&i.Slug,
		&i.Repo,
		pq.Array(&i.Keywords),
		&i.TeamkatalogenUrl,
	)
	return i, err
}
