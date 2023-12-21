package models

type DatasetComplete struct {
	Dataset
	Mappings   []MappingService `json:"mappings"`
	Access     []*Access        `json:"access"`
	Owner      *Owner           `json:"owner"`
	Services   *DatasetServices `json:"datasetServices"`
	Datasource Datasource       `json:"datasource"`
}

type BigQueryComplete struct {
	BigQuery
	Schema []*TableColumn `json:"schema"`
}

func (BigQueryComplete) IsDatasource() {}

type DataproductComplete struct {
	Dataproduct
	Datasets []*DatasetComplete `json:"datasets"`
	Keywords []string           `json:"keywords"`
}
