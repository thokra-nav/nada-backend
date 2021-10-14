// Code generated by sqlc. DO NOT EDIT.

package gensql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Datasource string

const (
	DatasourceBigquery Datasource = "bigquery"
)

func (e *Datasource) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Datasource(s)
	case string:
		*e = Datasource(s)
	default:
		return fmt.Errorf("unsupported scan type for Datasource: %T", src)
	}
	return nil
}

type Dataproduct struct {
	ID           uuid.UUID
	Name         string
	Description  sql.NullString
	Pii          bool
	Created      time.Time
	LastModified time.Time
	Type         Datasource
	TsvDocument  interface{}
}

type DataproductCollection struct {
	ID           uuid.UUID
	Name         string
	Description  sql.NullString
	Slug         string
	Repo         sql.NullString
	Created      time.Time
	LastModified time.Time
	Team         string
	Keywords     []string
	TsvDocument  interface{}
}

type DatasourceBigquery struct {
	DataproductID uuid.UUID
	ProjectID     string
	Dataset       string
	TableName     string
	Schema        json.RawMessage
}
