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

func (c *Bigquery) ComposePseudoViewQuery(projectID, datasetID, tableID string, targetColumns []string) string {
	qGenSalt := `WITH gen_salt AS (
		SELECT GENERATE_UUID() AS salt
	)`

	qSelect := "SELECT "
	for _, c := range targetColumns {
		qSelect += fmt.Sprintf(" SHA256(%v || gen_salt.salt) AS _x_%v", c, c)
		qSelect += ","
	}

	qSelect += "I.* EXCEPT("

	for i, c := range targetColumns {
		qSelect += c
		if i != len(targetColumns)-1 {
			qSelect += ","
		} else {
			qSelect += ")"
		}
	}
	qFrom := fmt.Sprintf("FROM `%v.%v.%v` AS I, gen_salt", projectID, datasetID, tableID)

	return qGenSalt + " " + qSelect + " " + qFrom
}

func (c *Bigquery) CreatePseudonymisedView(ctx context.Context, projectID, datasetID, tableID string, piiColumns []string) (string, string, string, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return "", "", "", fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	viewQuery := c.ComposePseudoViewQuery(projectID, datasetID, tableID, piiColumns)
	fmt.Println(viewQuery)
	meta := &bigquery.TableMetadata{
		ViewQuery: viewQuery,
	}
	pseudoViewID := fmt.Sprintf("_x_%v", tableID)
	if err := client.Dataset(datasetID).Table(pseudoViewID).Create(ctx, meta); err != nil {
		return "", "", "", err
	}
	return projectID, datasetID, pseudoViewID, nil
}

func (c *Bigquery) ComposeJoinableViewQuery(plainTableUrl models.BigQuery, joinableDatasetID string, pseudoColumns []string) string {
	qSalt := fmt.Sprintf("WITH unified_salt AS (SELECT value AS salt FROM `%v.%v.%v` ds WHERE ds.key='%v')", c.centralDataProject, "secrets_vault", "secrets", joinableDatasetID)

	qSelect := "SELECT "
	for _, c := range pseudoColumns {
		qSelect += fmt.Sprintf(" SHA256(%v || unified_salt.salt) AS _x_%v", c, c)
		qSelect += ","
	}

	qSelect += "I.* EXCEPT("

	for i, c := range pseudoColumns {
		qSelect += c
		if i != len(pseudoColumns)-1 {
			qSelect += ","
		} else {
			qSelect += ")"
		}
	}
	qFrom := fmt.Sprintf("FROM `%v.%v.%v` AS I, unified_salt", plainTableUrl.ProjectID, plainTableUrl.Dataset, plainTableUrl.Table)

	return qSalt + " " + qSelect + " " + qFrom
}

func (c *Bigquery) CreateJoinableView(ctx context.Context, joinableDatasetID string, tableUrl models.BigQuery) error {
	if !strings.HasPrefix(tableUrl.Table, "_x_") {
		return fmt.Errorf("invalid tableUrl: not a pseudo view")
	}

	plainTable := strings.TrimPrefix(tableUrl.Table, "_x_")

	client, err := bigquery.NewClient(ctx, tableUrl.ProjectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}

	meta, err := client.Dataset(tableUrl.Dataset).Table(tableUrl.Table).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("query table metadata: %v", err)
	}

	pseudoColumns := []string{}
	for _, field := range meta.Schema {
		if strings.HasPrefix(field.Name, "_x_") {
			pseudoColumns = append(pseudoColumns, strings.TrimPrefix(field.Name, "_x_"))
		}
	}

	if len(pseudoColumns) == 0 {
		return fmt.Errorf("invalid talbeUrl: no pseudo columns")
	}

	plainTableUrl := tableUrl
	plainTableUrl.Table = plainTable
	query := c.ComposeJoinableViewQuery(plainTableUrl, joinableDatasetID, pseudoColumns)

	centralProjectclient, err := bigquery.NewClient(ctx, c.centralDataProject)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer centralProjectclient.Close()

	joinableViewMeta := &bigquery.TableMetadata{
		ViewQuery: query,
	}

	if err := centralProjectclient.Dataset(joinableDatasetID).Table(tableUrl.Table).Create(ctx, joinableViewMeta); err != nil {
		return err
	}

	return nil
}

func (c *Bigquery) CreateJoinableViews(ctx context.Context, joinableDatasetID string, tableUrls []models.BigQuery) error {
	client, err := bigquery.NewClient(ctx, c.centralDataProject)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	c.createDatasetInCentralProject(ctx, joinableDatasetID)
	c.createSecretTable(ctx, "secrets_vault", "secrets")
	c.insertSecretIfNotExists(ctx, "secrets_vault", "secrets", joinableDatasetID)

	for _, table := range tableUrls {
		if err := c.CreateJoinableView(ctx, joinableDatasetID, table); err != nil {
			return err
		}
	}

	return nil
}

func (c *Bigquery) createDatasetInCentralProject(ctx context.Context, datasetID string) error {
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

	meta := &bigquery.DatasetMetadata{
		Location: "europe-north1",
	}

	if err := client.Dataset("secrets_vault").Create(ctx, meta); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code != 409 {
			return fmt.Errorf("failed to create secret dataset: %v", err)
		}
	}

	sampleSchema := bigquery.Schema{
		{Name: "key", Type: bigquery.StringFieldType},
		{Name: "value", Type: bigquery.StringFieldType},
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

func (c *Bigquery) insertSecretIfNotExists(ctx context.Context, secretDatasetID, secretTableID, key string) error {
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
	fmt.Fprintf(&insertQuery, "INSERT INTO `%v.%v.%v` (key, value) ", c.centralDataProject, secretDatasetID, secretTableID)
	fmt.Fprintf(&insertQuery, "SELECT '%v', '%v' FROM UNNEST([1]) ", key, encryptionKey.String())
	fmt.Fprintf(&insertQuery, "WHERE NOT EXISTS (SELECT 1 FROM `%v.%v.%v` WHERE key = '%v')", c.centralDataProject, secretDatasetID, secretTableID, key)

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
