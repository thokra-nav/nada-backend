package bigquery

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/graph/models"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

type Bigquery struct {
	centralDataProject string
}

func New(ctx context.Context, centralDataProject string) (*Bigquery, error) {
	return &Bigquery{
		centralDataProject: centralDataProject,
	}, nil
}

func (c *Bigquery) TableMetadata(ctx context.Context, projectID string, datasetID string, tableID string) (models.BigqueryMetadata, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return models.BigqueryMetadata{}, err
	}

	m, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return models.BigqueryMetadata{}, err
	}

	schema := models.BigquerySchema{}

	for _, c := range m.Schema {
		ct := "NULLABLE"
		switch {
		case c.Repeated:
			ct = "REPEATED"
		case c.Required:
			ct = "REQUIRED"
		}
		schema.Columns = append(schema.Columns, models.BigqueryColumn{
			Name:        c.Name,
			Type:        string(c.Type),
			Mode:        ct,
			Description: c.Description,
		})
	}

	metadata := models.BigqueryMetadata{
		Schema:       schema,
		LastModified: m.LastModifiedTime,
		Created:      m.CreationTime,
		Expires:      m.ExpirationTime,
		TableType:    m.Type,
		Description:  m.Description,
	}

	return metadata, nil
}

func (c *Bigquery) GetDatasets(ctx context.Context, projectID string) ([]string, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	datasets := []string{}
	it := client.Datasets(ctx)
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		datasets = append(datasets, ds.DatasetID)
	}
	return datasets, nil
}

func (c *Bigquery) GetTables(ctx context.Context, projectID, datasetID string) ([]*models.BigQueryTable, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	tables := []*models.BigQueryTable{}
	it := client.Dataset(datasetID).Tables(ctx)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}

		m, err := t.Metadata(ctx)
		if err != nil {
			return nil, err
		}

		if !isSupportedTableType(m.Type) {
			continue
		}

		tables = append(tables, &models.BigQueryTable{
			Name:         t.TableID,
			Description:  m.Description,
			Type:         models.BigQueryType(strings.ToLower(string(m.Type))),
			LastModified: m.LastModifiedTime,
		})
	}

	return tables, nil
}

func (c *Bigquery) ComposeViewQuery(secretDatasetID, secretTableID, encryptedDatasetId, encryptedTableId, projectID, datasetID, tableID string, piiColumns []string) string {
	var viewQuery strings.Builder
	fmt.Fprintf(&viewQuery, "WITH twts AS (SELECT ts.encryption_key AS encryption_key, * FROM `%v.%v.%v` t ", projectID, datasetID, tableID)
	fmt.Fprintf(&viewQuery, "LEFT JOIN `%v.%v.%v` ts ", c.centralDataProject, secretDatasetID, secretTableID)
	fmt.Fprintf(&viewQuery, "ON ts.table_id='%v.%v') SELECT ", encryptedDatasetId, encryptedTableId)

	exceptValues := " EXCEPT("
	for i, c := range piiColumns {
		if i != len(piiColumns)-1 {
			fmt.Fprintf(&viewQuery, "SHA256(%v) ^ SHA256('twts.encryption_key') AS x%v, ", c, c)
			exceptValues += fmt.Sprintf("%v,", c)
		} else {
			fmt.Fprintf(&viewQuery, "SHA256(%v) ^ SHA256('twts.encryption_key') AS x%v, * ", c, c)
			exceptValues += fmt.Sprintf("%v,encryption_key,table_id)", c)
		}
	}
	fmt.Fprintf(&viewQuery, "%v FROM twts", exceptValues)
	return viewQuery.String()
}

func (c *Bigquery) CreatePseudoynimizedView(ctx context.Context, projectID, datasetID, tableID string, piiColumns []string) (string, string, string, error) {
	client, err := bigquery.NewClient(ctx, c.centralDataProject)
	if err != nil {
		return "", "", "", fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	projectIDToUnderscore := strings.ReplaceAll(projectID, "-", "_")
	secretDatasetID := projectIDToUnderscore + "_vault"
	secretTableID := "secrets"
	encryptedDatasetId := "encrypted_views"
	encryptedTableId := fmt.Sprintf("%v_%v_%v", projectIDToUnderscore, datasetID, tableID)
	viewQuery := c.ComposeViewQuery(secretDatasetID, secretTableID, encryptedDatasetId, encryptedTableId, projectID, datasetID, tableID, piiColumns)
	if err := c.createSecretDataset(ctx, secretDatasetID); err != nil {
		return "", "", "", err
	}
	if err := c.createSecretTable(ctx, secretDatasetID, secretTableID); err != nil {
		return "", "", "", err
	}
	if err := c.insertEncryptionKeyIfNotExists(ctx, secretDatasetID, secretTableID, tableID); err != nil {
		return "", "", "", err
	}

	meta := &bigquery.TableMetadata{
		ViewQuery: viewQuery,
	}
	if err := client.Dataset(encryptedDatasetId).Table(encryptedTableId).Create(ctx, meta); err != nil {
		return "", "", "", err
	}
	return c.centralDataProject, encryptedDatasetId, encryptedTableId, nil
}

func (c *Bigquery) createSecretDataset(ctx context.Context, datasetID string) error {
	client, err := bigquery.NewClient(ctx, c.centralDataProject)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "europe-north1",
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 {
				return nil
			}
			return err
		}
	}
	return nil
}

func (c *Bigquery) createSecretTable(ctx context.Context, datasetID, tableID string) error {
	client, err := bigquery.NewClient(ctx, c.centralDataProject)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "table_id", Type: bigquery.StringFieldType},
		{Name: "encryption_key", Type: bigquery.StringFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 {
			return nil
		}
		return err
	}
	return nil
}

func (c *Bigquery) insertEncryptionKeyIfNotExists(ctx context.Context, secretDatasetID, secretTableID, tableID string) error {
	client, err := bigquery.NewClient(ctx, c.centralDataProject)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	encryptionKey, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	var insertQuery strings.Builder
	fmt.Fprintf(&insertQuery, "INSERT INTO `%v.%v.%v` (table_id, encryption_key) ", c.centralDataProject, secretDatasetID, secretTableID)
	fmt.Fprintf(&insertQuery, "SELECT '%v', '%v' FROM UNNEST([1]) ", tableID, encryptionKey.String())
	fmt.Fprintf(&insertQuery, "WHERE NOT EXISTS (SELECT 1 FROM `%v.%v.%v` WHERE table_id = '%v')", c.centralDataProject, secretDatasetID, secretTableID, tableID)

	job, err := client.Query(insertQuery.String()).Run(ctx)
	if err != nil {
		return err
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if status.Err() != nil {
		return err
	}

	return nil
}

func isSupportedTableType(tableType bigquery.TableType) bool {
	// We only support regular tables, views and materialized views for now.
	supported := []bigquery.TableType{
		bigquery.RegularTable,
		bigquery.ViewTable,
		bigquery.MaterializedView,
	}

	for _, tt := range supported {
		if tt == tableType {
			return true
		}
	}

	return false
}
