// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: tags.sql

package gensql

import (
	"context"
)

const createTagIfNotExist = `-- name: CreateTagIfNotExist :exec
INSERT INTO tags(phrase) VALUES ($1) ON CONFLICT DO NOTHING
`

func (q *Queries) CreateTagIfNotExist(ctx context.Context, phrase string) error {
	_, err := q.db.ExecContext(ctx, createTagIfNotExist, phrase)
	return err
}

const getKeywords = `-- name: GetKeywords :many
SELECT keyword::text, count(1) as "count"
FROM (
         SELECT unnest(ds.keywords) as keyword
            FROM datasets ds
         UNION ALL
         SELECT unnest(s.keywords) as keyword
            FROM stories s
    ) k
GROUP BY keyword
ORDER BY "count" DESC
`

type GetKeywordsRow struct {
	Keyword string
	Count   int64
}

func (q *Queries) GetKeywords(ctx context.Context) ([]GetKeywordsRow, error) {
	rows, err := q.db.QueryContext(ctx, getKeywords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetKeywordsRow{}
	for rows.Next() {
		var i GetKeywordsRow
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

const getTag = `-- name: GetTag :one
SELECT id, phrase FROM tags WHERE id=@id
`

func (q *Queries) GetTag(ctx context.Context) (Tag, error) {
	row := q.db.QueryRowContext(ctx, getTag)
	var i Tag
	err := row.Scan(&i.ID, &i.Phrase)
	return i, err
}

const getTagByPhrase = `-- name: GetTagByPhrase :one
SELECT id, phrase FROM tags WHERE phrase=@phrase
`

func (q *Queries) GetTagByPhrase(ctx context.Context) (Tag, error) {
	row := q.db.QueryRowContext(ctx, getTagByPhrase)
	var i Tag
	err := row.Scan(&i.ID, &i.Phrase)
	return i, err
}

const getTags = `-- name: GetTags :many
SELECT id, phrase FROM tags
`

func (q *Queries) GetTags(ctx context.Context) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, getTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Tag{}
	for rows.Next() {
		var i Tag
		if err := rows.Scan(&i.ID, &i.Phrase); err != nil {
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

const removeKeywordInDatasets = `-- name: RemoveKeywordInDatasets :exec
UPDATE datasets SET keywords= array_remove(keywords, $1)
`

func (q *Queries) RemoveKeywordInDatasets(ctx context.Context, keywordToRemove interface{}) error {
	_, err := q.db.ExecContext(ctx, removeKeywordInDatasets, keywordToRemove)
	return err
}

const removeKeywordInStories = `-- name: RemoveKeywordInStories :exec
UPDATE stories SET keywords= array_remove(keywords, $1)
`

func (q *Queries) RemoveKeywordInStories(ctx context.Context, keywordToRemove interface{}) error {
	_, err := q.db.ExecContext(ctx, removeKeywordInStories, keywordToRemove)
	return err
}

const replaceKeywordInDatasets = `-- name: ReplaceKeywordInDatasets :exec
UPDATE datasets SET keywords= array_replace(keywords, $1, $2)
`

type ReplaceKeywordInDatasetsParams struct {
	Keyword           interface{}
	NewTextForKeyword interface{}
}

func (q *Queries) ReplaceKeywordInDatasets(ctx context.Context, arg ReplaceKeywordInDatasetsParams) error {
	_, err := q.db.ExecContext(ctx, replaceKeywordInDatasets, arg.Keyword, arg.NewTextForKeyword)
	return err
}

const replaceKeywordInStories = `-- name: ReplaceKeywordInStories :exec
UPDATE stories SET keywords= array_replace(keywords, $1, $2)
`

type ReplaceKeywordInStoriesParams struct {
	Keyword           interface{}
	NewTextForKeyword interface{}
}

func (q *Queries) ReplaceKeywordInStories(ctx context.Context, arg ReplaceKeywordInStoriesParams) error {
	_, err := q.db.ExecContext(ctx, replaceKeywordInStories, arg.Keyword, arg.NewTextForKeyword)
	return err
}

const updateTag = `-- name: UpdateTag :exec
UPDATE tags SET phrase = $1 where phrase = $2
`

type UpdateTagParams struct {
	NewPhrase string
	OldPhrase string
}

func (q *Queries) UpdateTag(ctx context.Context, arg UpdateTagParams) error {
	_, err := q.db.ExecContext(ctx, updateTag, arg.NewPhrase, arg.OldPhrase)
	return err
}
