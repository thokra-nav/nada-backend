// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: team.sql

package gensql

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const deleteNadaToken = `-- name: DeleteNadaToken :exec
DELETE FROM
    nada_tokens
WHERE
    team = $1
`

func (q *Queries) DeleteNadaToken(ctx context.Context, team string) error {
	_, err := q.db.ExecContext(ctx, deleteNadaToken, team)
	return err
}

const getNadaToken = `-- name: GetNadaToken :one
SELECT
    token
FROM
    nada_tokens
WHERE
    team = $1
`

func (q *Queries) GetNadaToken(ctx context.Context, team string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getNadaToken, team)
	var token uuid.UUID
	err := row.Scan(&token)
	return token, err
}

const getNadaTokens = `-- name: GetNadaTokens :many
SELECT
    team, token
FROM
    nada_tokens
WHERE
    team = ANY ($1 :: text [])
ORDER BY
    team
`

func (q *Queries) GetNadaTokens(ctx context.Context, teams []string) ([]NadaToken, error) {
	rows, err := q.db.QueryContext(ctx, getNadaTokens, pq.Array(teams))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []NadaToken{}
	for rows.Next() {
		var i NadaToken
		if err := rows.Scan(&i.Team, &i.Token); err != nil {
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
