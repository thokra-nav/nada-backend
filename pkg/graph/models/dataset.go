package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/database/gensql"
)

type Dataset struct {
	ID                       uuid.UUID `json:"id"`
	DataproductID            uuid.UUID `json:"dataproductID"`
	Name                     string    `json:"name"`
	Created                  time.Time `json:"created"`
	LastModified             time.Time `json:"lastModified"`
	Description              *string   `json:"description"`
	Slug                     string    `json:"slug"`
	Repo                     *string   `json:"repo"`
	Pii                      PiiLevel  `json:"pii"`
	Keywords                 []string  `json:"keywords"`
	Type                     gensql.DatasourceType
	AnonymisationDescription *string `json:"anonymisationDescription"`
	TargetUser               *string `json:"targetUser"`
}

func (Dataset) IsSearchResult() {}

type Datasource interface {
	IsDatasource()
}

type BigQuery struct {
	DatasetID    uuid.UUID
	ProjectID    string       `json:"projectID"`
	Dataset      string       `json:"dataset"`
	Table        string       `json:"table"`
	TableType    BigQueryType `json:"tableType"`
	LastModified time.Time    `json:"lastModified"`
	Created      time.Time    `json:"created"`
	Expires      *time.Time   `json:"expired"`
	Description  string       `json:"description"`
	PiiTags      *string      `json:"piiTags"`
	MissingSince *time.Time   `json:"missingSince"`
}

func (BigQuery) IsDatasource() {}

type NewBigQuery struct {
	ProjectID string  `json:"projectID"`
	Dataset   string  `json:"dataset"`
	Table     string  `json:"table"`
	PiiTags   *string `json:"piiTags"`
}

type NewDataset struct {
	DataproductID            uuid.UUID   `json:"dataproductID"`
	Name                     string      `json:"name"`
	Description              *string     `json:"description"`
	Slug                     *string     `json:"slug"`
	Repo                     *string     `json:"repo"`
	Pii                      PiiLevel    `json:"pii"`
	Keywords                 []string    `json:"keywords"`
	BigQuery                 NewBigQuery `json:"bigquery"`
	AnonymisationDescription *string     `json:"anonymisationDescription"`
	GrantAllUsers            *bool       `json:"grantAllUsers"`
	TargetUser               *string     `json:"targetUser"`
	CreatePseudoynimizedView bool        `json:"createPseudoynimizedView"`
	Metadata                 BigqueryMetadata
}

// NewDatasetForNewDataproduct contains metadata for creating a new dataset for a new dataproduct
type NewDatasetForNewDataproduct struct {
	Name                     string      `json:"name"`
	Description              *string     `json:"description"`
	Repo                     *string     `json:"repo"`
	Pii                      PiiLevel    `json:"pii"`
	Keywords                 []string    `json:"keywords"`
	Bigquery                 NewBigQuery `json:"bigquery"`
	AnonymisationDescription *string     `json:"anonymisationDescription"`
	GrantAllUsers            *bool       `json:"grantAllUsers"`
	TargetUser               *string     `json:"targetUser"`
	CreatePseudoynimizedView bool        `json:"createPseudoynimizedView"`
	Metadata                 BigqueryMetadata
}

type UpdateDataset struct {
	Name                     string     `json:"name"`
	Description              *string    `json:"description"`
	Slug                     *string    `json:"slug"`
	Repo                     *string    `json:"repo"`
	Pii                      PiiLevel   `json:"pii"`
	Keywords                 []string   `json:"keywords"`
	DataproductID            *uuid.UUID `json:"dataproductID"`
	AnonymisationDescription *string    `json:"anonymisationDescription"`
	PiiTags                  *string    `json:"piiTags"`
	TargetUser               *string    `json:"targetUser"`
	CreatePseudoynimizedView bool       `json:"createPseudoynimizedView"`
}

type DatasetServices struct {
	Metabase *string `json:"metabase"`
}

func ToNewDataset(dpID uuid.UUID, dsForDP *NewDatasetForNewDataproduct) *NewDataset {
	if dsForDP.Keywords == nil {
		dsForDP.Keywords = []string{}
	}
	return &NewDataset{
		DataproductID:            dpID,
		Name:                     dsForDP.Name,
		Description:              dsForDP.Description,
		Repo:                     dsForDP.Repo,
		Pii:                      dsForDP.Pii,
		Keywords:                 dsForDP.Keywords,
		BigQuery:                 dsForDP.Bigquery,
		AnonymisationDescription: dsForDP.AnonymisationDescription,
		GrantAllUsers:            dsForDP.GrantAllUsers,
		TargetUser:               dsForDP.TargetUser,
		CreatePseudoynimizedView: dsForDP.CreatePseudoynimizedView,
		Metadata:                 dsForDP.Metadata,
	}
}

func ToNewDatasetForNewDataproduct(ds *NewDataset) *NewDatasetForNewDataproduct {
	if ds.Keywords == nil {
		ds.Keywords = []string{}
	}

	return &NewDatasetForNewDataproduct{
		Name:                     ds.Name,
		Description:              ds.Description,
		Repo:                     ds.Repo,
		Pii:                      ds.Pii,
		Keywords:                 ds.Keywords,
		Bigquery:                 ds.BigQuery,
		AnonymisationDescription: ds.AnonymisationDescription,
		GrantAllUsers:            ds.GrantAllUsers,
		TargetUser:               ds.TargetUser,
		CreatePseudoynimizedView: ds.CreatePseudoynimizedView,
		Metadata:                 ds.Metadata,
	}
}
