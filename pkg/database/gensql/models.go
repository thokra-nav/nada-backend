// Code generated by sqlc. DO NOT EDIT.

package gensql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tabbed/pqtype"
)

type AccessRequestStatusType string

const (
	AccessRequestStatusTypePending  AccessRequestStatusType = "pending"
	AccessRequestStatusTypeApproved AccessRequestStatusType = "approved"
	AccessRequestStatusTypeDenied   AccessRequestStatusType = "denied"
)

func (e *AccessRequestStatusType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccessRequestStatusType(s)
	case string:
		*e = AccessRequestStatusType(s)
	default:
		return fmt.Errorf("unsupported scan type for AccessRequestStatusType: %T", src)
	}
	return nil
}

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

type StoryViewType string

const (
	StoryViewTypeMarkdown StoryViewType = "markdown"
	StoryViewTypeHeader   StoryViewType = "header"
	StoryViewTypePlotly   StoryViewType = "plotly"
	StoryViewTypeVega     StoryViewType = "vega"
)

func (e *StoryViewType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StoryViewType(s)
	case string:
		*e = StoryViewType(s)
	default:
		return fmt.Errorf("unsupported scan type for StoryViewType: %T", src)
	}
	return nil
}

type Dataproduct struct {
	ID               uuid.UUID
	Name             string
	Description      sql.NullString
	Group            string
	Pii              bool
	Created          time.Time
	LastModified     time.Time
	Type             DatasourceType
	TsvDocument      interface{}
	Slug             string
	Repo             sql.NullString
	Keywords         []string
	TeamkatalogenUrl sql.NullString
}

type DataproductAccess struct {
	ID              uuid.UUID
	DataproductID   uuid.UUID
	Subject         string
	Granter         string
	Expires         sql.NullTime
	Created         time.Time
	Revoked         sql.NullTime
	AccessRequestID uuid.NullUUID
}

type DataproductAccessRequest struct {
	ID                   uuid.UUID
	DataproductID        uuid.UUID
	Subject              string
	Owner                string
	PollyDocumentationID uuid.NullUUID
	LastModified         time.Time
	Created              time.Time
	Expires              sql.NullTime
	Status               AccessRequestStatusType
	Closed               sql.NullTime
	Granter              sql.NullString
	Reason               sql.NullString
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
	Description   sql.NullString
}

type MetabaseMetadatum struct {
	DataproductID     uuid.UUID
	DatabaseID        int32
	PermissionGroupID sql.NullInt32
	SaEmail           string
	CollectionID      sql.NullInt32
	DeletedAt         sql.NullTime
}

type PollyDocumentation struct {
	ID         uuid.UUID
	ExternalID string
	Name       string
	Url        string
}

type Search struct {
	ElementID    uuid.UUID
	ElementType  interface{}
	Description  string
	Keywords     []string
	Group        string
	Created      time.Time
	LastModified time.Time
	TsvDocument  interface{}
	Services     []string
}

type Session struct {
	Token       string
	AccessToken string
	Email       string
	Name        string
	Created     time.Time
	Expires     time.Time
}

type Story struct {
	ID               uuid.UUID
	Name             string
	Created          time.Time
	LastModified     time.Time
	Group            string
	Description      sql.NullString
	Keywords         []string
	TeamkatalogenUrl sql.NullString
}

type StoryDraft struct {
	ID      uuid.UUID
	Name    string
	Created time.Time
}

type StoryToken struct {
	ID      uuid.UUID
	StoryID uuid.UUID
	Token   uuid.UUID
}

type StoryView struct {
	ID      uuid.UUID
	StoryID uuid.UUID
	Sort    int32
	Type    StoryViewType
	Spec    json.RawMessage
}

type StoryViewDraft struct {
	ID      uuid.UUID
	StoryID uuid.UUID
	Sort    int32
	Type    StoryViewType
	Spec    json.RawMessage
}

type ThirdPartyMapping struct {
	DataproductID uuid.UUID
	Services      []string
}
