package gensql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// LoggingDBTX struct implementing the DBTX interface as a delegate
type LoggingDBTX struct {
	db DBTX
}

// New function modified to return a Queries object with db set to the delegate of the parameter db
func NewWithLog(db DBTX) *Queries {
	return &Queries{db: &LoggingDBTX{db: db}}
}

// WithTx function in Queries struct
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return q
}

// Logging methods for LoggingDBTX
func (l *LoggingDBTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	queryComponents := strings.Split(query, ":")
	fmt.Printf("ExecContext called with query: %s and args: %v\n", queryComponents[1], args)
	return l.db.ExecContext(ctx, query, args...)
}

func (l *LoggingDBTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	queryComponents := strings.Split(query, ":")
	fmt.Printf("PrepareContext called with query: %s\n", queryComponents[1])
	return l.db.PrepareContext(ctx, query)
}

func (l *LoggingDBTX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	queryComponents := strings.Split(query, ":")
	fmt.Printf("QueryContext called with query: %s and args: %v\n", queryComponents[1], args)
	return l.db.QueryContext(ctx, query, args...)
}

func (l *LoggingDBTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	queryComponents := strings.Split(query, ":")
	fmt.Printf("QueryRowContext called with query: %s and args: %v\n", queryComponents[1], args)
	return l.db.QueryRowContext(ctx, query, args...)
}
