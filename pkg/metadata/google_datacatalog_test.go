package metadata

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/navikt/nada-backend/pkg/openapi"
	datacatalogpb "google.golang.org/genproto/googleapis/cloud/datacatalog/v1"
)

func TestGoogleDatacatalog(t *testing.T) {
	expected := []openapi.BigQuery{
		{
			ProjectId: "project_id",
			Dataset:   "mydataset",
			Table:     "mytable",
		},
	}
	dcc := &googleMockClient{
		searchResponse: []*datacatalogpb.SearchCatalogResult{
			{LinkedResource: "//some/ completely /// other / type / of / resource"},
			{LinkedResource: "//bigquery.googleapis.com/projects/project_id/datasets/mydataset/tables/mytable"},
			{LinkedResource: "//bigquery.googleapis.com/projects/project_id/datasets/mydataset"},
		},
	}
	client := &Datacatalog{client: dcc}

	res, err := client.GetDatasets(context.Background(), "project_id")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(res, expected) {
		t.Error(cmp.Diff(res, expected))
	}
}

type googleMockClient struct {
	searchResponse []*datacatalogpb.SearchCatalogResult
	lookupResponse *datacatalogpb.Entry
	err            error
}

func (g *googleMockClient) SearchCatalog(ctx context.Context, req *datacatalogpb.SearchCatalogRequest) ([]*datacatalogpb.SearchCatalogResult, error) {
	return g.searchResponse, g.err
}

func (g *googleMockClient) LookupEntry(ctx context.Context, req *datacatalogpb.LookupEntryRequest) (*datacatalogpb.Entry, error) {
	return g.lookupResponse, g.err
}

func (g *googleMockClient) Close() error { return nil }
