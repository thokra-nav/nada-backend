// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package gensql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
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

type NullAccessRequestStatusType struct {
	AccessRequestStatusType AccessRequestStatusType
	Valid                   bool // Valid is true if AccessRequestStatusType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccessRequestStatusType) Scan(value interface{}) error {
	if value == nil {
		ns.AccessRequestStatusType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccessRequestStatusType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccessRequestStatusType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccessRequestStatusType), nil
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

type NullDatasourceType struct {
	DatasourceType DatasourceType
	Valid          bool // Valid is true if DatasourceType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullDatasourceType) Scan(value interface{}) error {
	if value == nil {
		ns.DatasourceType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.DatasourceType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullDatasourceType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.DatasourceType), nil
}

type PiiLevel string

const (
	PiiLevelSensitive  PiiLevel = "sensitive"
	PiiLevelAnonymised PiiLevel = "anonymised"
	PiiLevelNone       PiiLevel = "none"
)

func (e *PiiLevel) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PiiLevel(s)
	case string:
		*e = PiiLevel(s)
	default:
		return fmt.Errorf("unsupported scan type for PiiLevel: %T", src)
	}
	return nil
}

type NullPiiLevel struct {
	PiiLevel PiiLevel
	Valid    bool // Valid is true if PiiLevel is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPiiLevel) Scan(value interface{}) error {
	if value == nil {
		ns.PiiLevel, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PiiLevel.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPiiLevel) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PiiLevel), nil
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

type NullStoryViewType struct {
	StoryViewType StoryViewType
	Valid         bool // Valid is true if StoryViewType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStoryViewType) Scan(value interface{}) error {
	if value == nil {
		ns.StoryViewType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StoryViewType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStoryViewType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StoryViewType), nil
}

type Dashboard struct {
	ID  string
	Url string
}

type Dataproduct struct {
	ID               uuid.UUID
	Name             string
	Description      sql.NullString
	Group            string
	Created          time.Time
	LastModified     time.Time
	TsvDocument      interface{}
	Slug             string
	TeamkatalogenUrl sql.NullString
	TeamContact      sql.NullString
	TeamID           sql.NullString
}

type DataproductCompleteView struct {
	DataproductID    uuid.UUID
	DpName           string
	DpDescription    sql.NullString
	DpGroup          string
	DpCreated        time.Time
	DpLastModified   time.Time
	DpSlug           string
	TeamkatalogenUrl sql.NullString
	TeamContact      sql.NullString
	TeamID           sql.NullString
	BqID             uuid.UUID
	BqCreated        time.Time
	BqLastModified   time.Time
	BqExpires        sql.NullTime
	BqDescription    sql.NullString
	BqMissingSince   sql.NullTime
	PiiTags          pqtype.NullRawMessage
	BqProject        string
	BqDataset        string
	BqTableName      string
	BqTableType      string
	PseudoColumns    []string
	BqSchema         pqtype.NullRawMessage
	DsDpID           uuid.NullUUID
	DsID             uuid.NullUUID
	DsName           sql.NullString
	DsDescription    sql.NullString
	DsCreated        sql.NullTime
	DsLastModified   sql.NullTime
	DsSlug           sql.NullString
	DsKeywords       []string
	MappingServices  []string
	AccessID         uuid.NullUUID
	AccessSubject    sql.NullString
	AccessGranter    sql.NullString
	AccessExpires    sql.NullTime
	AccessCreated    sql.NullTime
	AccessRevoked    sql.NullTime
	AccessRequestID  uuid.NullUUID
	MbDatabaseID     sql.NullInt32
}

type Dataset struct {
	ID                       uuid.UUID
	Name                     string
	Description              sql.NullString
	Pii                      PiiLevel
	Created                  time.Time
	LastModified             time.Time
	Type                     DatasourceType
	TsvDocument              interface{}
	Slug                     string
	Repo                     sql.NullString
	Keywords                 []string
	DataproductID            uuid.UUID
	AnonymisationDescription sql.NullString
	TargetUser               sql.NullString
}

type DatasetAccess struct {
	ID              uuid.UUID
	DatasetID       uuid.UUID
	Subject         string
	Granter         string
	Expires         sql.NullTime
	Created         time.Time
	Revoked         sql.NullTime
	AccessRequestID uuid.NullUUID
}

type DatasetAccessRequest struct {
	ID                   uuid.UUID
	DatasetID            uuid.UUID
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

type DatasourceBigquery struct {
	DatasetID     uuid.UUID
	ProjectID     string
	Dataset       string
	TableName     string
	Schema        pqtype.NullRawMessage
	LastModified  time.Time
	Created       time.Time
	Expires       sql.NullTime
	TableType     string
	Description   sql.NullString
	PiiTags       pqtype.NullRawMessage
	MissingSince  sql.NullTime
	ID            uuid.UUID
	IsReference   bool
	PseudoColumns []string
	Deleted       sql.NullTime
}

type HttpCache struct {
	ID                int32
	Endpoint          string
	ResponseBody      []byte
	CreatedAt         time.Time
	LastTriedUpdateAt time.Time
}

type InsightProduct struct {
	ID               uuid.UUID
	Name             string
	Description      sql.NullString
	Creator          string
	Created          time.Time
	LastModified     time.Time
	Type             string
	TsvDocument      interface{}
	Link             string
	Keywords         []string
	Group            string
	TeamkatalogenUrl sql.NullString
	TeamID           sql.NullString
}

type JoinableView struct {
	ID      uuid.UUID
	Owner   string
	Name    string
	Created time.Time
	Expires sql.NullTime
	Deleted sql.NullTime
}

type JoinableViewsDatasource struct {
	ID             uuid.UUID
	JoinableViewID uuid.UUID
	DatasourceID   uuid.UUID
	Deleted        sql.NullTime
}

type MetabaseMetadatum struct {
	DatabaseID        int32
	PermissionGroupID sql.NullInt32
	SaEmail           string
	CollectionID      sql.NullInt32
	DeletedAt         sql.NullTime
	DatasetID         uuid.UUID
}

type NadaToken struct {
	Team  string
	Token uuid.UUID
}

type PollyDocumentation struct {
	ID         uuid.UUID
	ExternalID string
	Name       string
	Url        string
}

type QuartoStory struct {
	ID               uuid.UUID
	Name             string
	Creator          string
	Created          time.Time
	LastModified     time.Time
	Description      string
	Keywords         []string
	TeamkatalogenUrl sql.NullString
	TeamID           sql.NullString
	Group            string
}

type Search struct {
	ElementID    uuid.UUID
	ElementType  string
	Description  string
	Keywords     interface{}
	Group        string
	TeamID       sql.NullString
	Created      time.Time
	LastModified time.Time
	TsvDocument  interface{}
	Services     string
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
	TeamID           sql.NullString
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

type Tag struct {
	ID     uuid.UUID
	Phrase string
}

type TeamProject struct {
	Team    string
	Project string
}

type ThirdPartyMapping struct {
	Services  []string
	DatasetID uuid.UUID
}
