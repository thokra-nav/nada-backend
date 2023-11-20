// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: dashboards.sql

package gensql

import (
	"context"
)

const getDashboard = `-- name: GetDashboard :one
SELECT id, url
FROM "dashboards"
WHERE id = $1
`

func (q *Queries) GetDashboard(ctx context.Context, id string) (Dashboard, error) {
	row := q.db.QueryRowContext(ctx, getDashboard, id)
	var i Dashboard
	err := row.Scan(&i.ID, &i.Url)
	return i, err
}
