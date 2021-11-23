// Code generated by sqlc. DO NOT EDIT.
// source: collections.sql

package gensql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createCollection = `-- name: CreateCollection :one
INSERT INTO collections (
	"name",
	"description",
	"slug",
	"group",
	"teamkatalogen_url",
	"keywords"
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
) RETURNING id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url
`

type CreateCollectionParams struct {
	Name                  string
	Description           sql.NullString
	Slug                  string
	OwnerGroup            string
	OwnerTeamkatalogenUrl sql.NullString
	Keywords              []string
}

func (q *Queries) CreateCollection(ctx context.Context, arg CreateCollectionParams) (Collection, error) {
	row := q.db.QueryRowContext(ctx, createCollection,
		arg.Name,
		arg.Description,
		arg.Slug,
		arg.OwnerGroup,
		arg.OwnerTeamkatalogenUrl,
		pq.Array(arg.Keywords),
	)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Slug,
		&i.Created,
		&i.LastModified,
		&i.Group,
		pq.Array(&i.Keywords),
		&i.TsvDocument,
		&i.TeamkatalogenUrl,
	)
	return i, err
}

const createCollectionElement = `-- name: CreateCollectionElement :exec
INSERT INTO collection_elements (
	"element_id",
	"collection_id",
	"element_type"
) VALUES (
	$1,
	$2,
	$3
)
`

type CreateCollectionElementParams struct {
	ElementID    uuid.UUID
	CollectionID uuid.UUID
	ElementType  string
}

func (q *Queries) CreateCollectionElement(ctx context.Context, arg CreateCollectionElementParams) error {
	_, err := q.db.ExecContext(ctx, createCollectionElement, arg.ElementID, arg.CollectionID, arg.ElementType)
	return err
}

const deleteCollection = `-- name: DeleteCollection :exec
DELETE FROM collections WHERE id = $1
`

func (q *Queries) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteCollection, id)
	return err
}

const deleteCollectionElement = `-- name: DeleteCollectionElement :exec
DELETE FROM collection_elements WHERE element_id = $1 AND collection_id = $2 AND element_type = $3
`

type DeleteCollectionElementParams struct {
	ElementID    uuid.UUID
	CollectionID uuid.UUID
	ElementType  string
}

func (q *Queries) DeleteCollectionElement(ctx context.Context, arg DeleteCollectionElementParams) error {
	_, err := q.db.ExecContext(ctx, deleteCollectionElement, arg.ElementID, arg.CollectionID, arg.ElementType)
	return err
}

const getCollection = `-- name: GetCollection :one
SELECT id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url FROM collections WHERE id = $1
`

func (q *Queries) GetCollection(ctx context.Context, id uuid.UUID) (Collection, error) {
	row := q.db.QueryRowContext(ctx, getCollection, id)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Slug,
		&i.Created,
		&i.LastModified,
		&i.Group,
		pq.Array(&i.Keywords),
		&i.TsvDocument,
		&i.TeamkatalogenUrl,
	)
	return i, err
}

const getCollectionElements = `-- name: GetCollectionElements :many
SELECT id, name, description, "group", pii, created, last_modified, type, tsv_document, slug, repo, keywords, teamkatalogen_url
FROM dataproducts
WHERE id IN
	(SELECT element_id FROM collection_elements WHERE collection_id = $1 AND element_type = 'dataproduct')
`

func (q *Queries) GetCollectionElements(ctx context.Context, collectionID uuid.UUID) ([]Dataproduct, error) {
	rows, err := q.db.QueryContext(ctx, getCollectionElements, collectionID)
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

const getCollections = `-- name: GetCollections :many
SELECT id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url FROM collections ORDER BY last_modified DESC LIMIT $2 OFFSET $1
`

type GetCollectionsParams struct {
	Offset int32
	Limit  int32
}

func (q *Queries) GetCollections(ctx context.Context, arg GetCollectionsParams) ([]Collection, error) {
	rows, err := q.db.QueryContext(ctx, getCollections, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Collection{}
	for rows.Next() {
		var i Collection
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Slug,
			&i.Created,
			&i.LastModified,
			&i.Group,
			pq.Array(&i.Keywords),
			&i.TsvDocument,
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

const getCollectionsByGroups = `-- name: GetCollectionsByGroups :many
SELECT id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url FROM collections WHERE "group" = ANY($1::text[]) ORDER BY last_modified DESC
`

func (q *Queries) GetCollectionsByGroups(ctx context.Context, groups []string) ([]Collection, error) {
	rows, err := q.db.QueryContext(ctx, getCollectionsByGroups, pq.Array(groups))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Collection{}
	for rows.Next() {
		var i Collection
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Slug,
			&i.Created,
			&i.LastModified,
			&i.Group,
			pq.Array(&i.Keywords),
			&i.TsvDocument,
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

const getCollectionsByIDs = `-- name: GetCollectionsByIDs :many
SELECT id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url FROM collections WHERE id = ANY($1::uuid[]) ORDER BY last_modified DESC
`

func (q *Queries) GetCollectionsByIDs(ctx context.Context, ids []uuid.UUID) ([]Collection, error) {
	rows, err := q.db.QueryContext(ctx, getCollectionsByIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Collection{}
	for rows.Next() {
		var i Collection
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Slug,
			&i.Created,
			&i.LastModified,
			&i.Group,
			pq.Array(&i.Keywords),
			&i.TsvDocument,
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

const getCollectionsForElement = `-- name: GetCollectionsForElement :many
SELECT id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url
FROM collections
WHERE id IN
	(SELECT collection_id FROM collection_elements WHERE element_id = $1 AND element_type = $2)
LIMIT $4
OFFSET $3
`

type GetCollectionsForElementParams struct {
	ElementID   uuid.UUID
	ElementType string
	Offset      int32
	Limit       int32
}

func (q *Queries) GetCollectionsForElement(ctx context.Context, arg GetCollectionsForElementParams) ([]Collection, error) {
	rows, err := q.db.QueryContext(ctx, getCollectionsForElement,
		arg.ElementID,
		arg.ElementType,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Collection{}
	for rows.Next() {
		var i Collection
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Slug,
			&i.Created,
			&i.LastModified,
			&i.Group,
			pq.Array(&i.Keywords),
			&i.TsvDocument,
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

const updateCollection = `-- name: UpdateCollection :one
UPDATE collections SET
	"name" = $1,
	"description" = $2,
	"slug" = $3,
	"teamkatalogen_url" = $4,
	"keywords" = $5
WHERE id = $6
RETURNING id, name, description, slug, created, last_modified, "group", keywords, tsv_document, teamkatalogen_url
`

type UpdateCollectionParams struct {
	Name                  string
	Description           sql.NullString
	Slug                  string
	OwnerTeamkatalogenUrl sql.NullString
	Keywords              []string
	ID                    uuid.UUID
}

func (q *Queries) UpdateCollection(ctx context.Context, arg UpdateCollectionParams) (Collection, error) {
	row := q.db.QueryRowContext(ctx, updateCollection,
		arg.Name,
		arg.Description,
		arg.Slug,
		arg.OwnerTeamkatalogenUrl,
		pq.Array(arg.Keywords),
		arg.ID,
	)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Slug,
		&i.Created,
		&i.LastModified,
		&i.Group,
		pq.Array(&i.Keywords),
		&i.TsvDocument,
		&i.TeamkatalogenUrl,
	)
	return i, err
}
