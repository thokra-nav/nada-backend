// Code generated by sqlc. DO NOT EDIT.

package gensql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tabbed/pqtype"
)

type DatasourceType string

const (
	DatasourceTypeBigquery DatasourceType = "bigquery"
)

func (e *DatasourceType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = DatasourceType(s)
	case string:
		*e = DatasourceType(s)
	default:
		return fmt.Errorf("unsupported scan type for DatasourceType: %T", src)
	}
	return nil
}

type Collection struct {
	ID           uuid.UUID
	Name         string
	Description  sql.NullString
	Slug         string
	Created      time.Time
	LastModified time.Time
	Group        string
	Keywords     []string
	TsvDocument  interface{}
}

type CollectionElement struct {
	ElementID    uuid.UUID
	CollectionID uuid.UUID
	ElementType  string
}

type Dataproduct struct {
	ID           uuid.UUID
	Name         string
	Description  sql.NullString
	Group        string
	Pii          bool
	Created      time.Time
	LastModified time.Time
	Type         DatasourceType
	TsvDocument  interface{}
	Slug         string
	Repo         sql.NullString
	Keywords     []string
}

type DataproductAccess struct {
	ID            uuid.UUID
	DataproductID uuid.UUID
	Subject       string
	Granter       string
	Expires       sql.NullTime
	Created       time.Time
	Revoked       sql.NullTime
}

type DataproductRequester struct {
	DataproductID uuid.UUID
	Subject       string
}

type DatasourceBigquery struct {
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

type Search struct {
	ElementID    uuid.UUID
	ElementType  interface{}
	LastModified time.Time
	Keywords     []string
	Group        string
	Created      time.Time
	TsvDocument  interface{}
}
