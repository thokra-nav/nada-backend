package graph

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/99designs/gqlgen-contrib/prometheus"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/auth"
	bq "github.com/navikt/nada-backend/pkg/bigquery"
	"github.com/navikt/nada-backend/pkg/database"
	"github.com/navikt/nada-backend/pkg/graph/generated"
	"github.com/navikt/nada-backend/pkg/graph/models"
	"github.com/sirupsen/logrus"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type Bigquery interface {
	GetTables(ctx context.Context, projectID, datasetID string) ([]*models.BigQueryTable, error)
	GetDatasets(ctx context.Context, projectID string) ([]string, error)
	TableMetadata(ctx context.Context, projectID string, datasetID string, tableID string) (models.BigqueryMetadata, error)
	CreatePseudonymisedView(ctx context.Context, projectID string, datasetID string, tableID string, targetColumns []string) (string, string, string, error)
	CreateJoinableViewsForUser(ctx context.Context, name string, datasources []bq.JoinableViewDatasource) (string, string, map[uuid.UUID]string, error)
	MakeBigQueryUrlForJoinableViews(name, projectID, datasetID, tableID string) string
	DeleteJoinableDataset(ctx context.Context, datasetID string) error
	DeleteJoinableView(ctx context.Context, joinableViewName, refProjectID, refDatasetID, refTableID string) error
	DeletePseudoView(ctx context.Context, pseudoProjectID, pseudoDatasetID, pseudoTableID string) error
}

type AccessManager interface {
	Grant(ctx context.Context, projectID, dataset, table, member string) error
	Revoke(ctx context.Context, projectID, dataset, table, member string) error
	AddToAuthorizedViews(ctx context.Context, srcProjectID, srcDataset, sinkProjectID, sinkDataset, sinkTable string) error
}

type Polly interface {
	SearchPolly(ctx context.Context, q string) ([]*models.QueryPolly, error)
}

type Teamkatalogen interface {
	Search(ctx context.Context, query string) ([]*models.TeamkatalogenResult, error)
	GetTeamsInProductArea(ctx context.Context, paID string) ([]*models.Team, error)
	GetProductArea(ctx context.Context, paID string) (*models.ProductArea, error)
	GetProductAreas(ctx context.Context) ([]*models.ProductArea, error)
	GetTeam(ctx context.Context, teamID string) (*models.Team, error)
	GetTeamCatalogURL(teamID string) string
}

type Slack interface {
	NewAccessRequest(contact string, dp *models.Dataproduct, ds *models.Dataset, ar *models.AccessRequest) error
	IsValidSlackChannel(name string) (bool, error)
}

type Resolver struct {
	repo               *database.Repo
	bigquery           Bigquery
	gcpProjects        *auth.TeamProjectsMapping
	accessMgr          AccessManager
	teamkatalogen      Teamkatalogen
	slack              Slack
	pollyAPI           Polly
	centralDataProject string
	log                *logrus.Entry
}

func New(repo *database.Repo, gcp Bigquery, gcpProjects *auth.TeamProjectsMapping, accessMgr AccessManager, tk Teamkatalogen, slack Slack, pollyAPI Polly, centralDataProject string, log *logrus.Entry) *handler.Server {
	resolver := &Resolver{
		repo:               repo,
		bigquery:           gcp,
		gcpProjects:        gcpProjects,
		accessMgr:          accessMgr,
		teamkatalogen:      tk,
		slack:              slack,
		pollyAPI:           pollyAPI,
		centralDataProject: centralDataProject,
		log:                log,
	}

	config := generated.Config{Resolvers: resolver}
	config.Directives.Authenticated = authenticate
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	srv.Use(prometheus.Tracer{})
	return srv
}

func (r *Resolver) ensureGroupOwnsGCPProject(ctx context.Context, group, projectID string) error {
	groupProject, ok := r.gcpProjects.Get(group)
	if !ok {
		return ErrUnauthorized
	}

	if groupProject == projectID {
		return nil
	}

	return ErrUnauthorized
}

func pagination(limit *int, offset *int) (int, int) {
	l := 15
	o := 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}
	return l, o
}

func ensureOwner(ctx context.Context, owner string) error {
	user := auth.GetUser(ctx)

	if user != nil && (user.GoogleGroups.Contains(owner) || owner == user.Email) {
		return nil
	}

	return ErrUnauthorized
}

func ensureUserInGroup(ctx context.Context, group string) error {
	user := auth.GetUser(ctx)
	if user == nil || !user.GoogleGroups.Contains(group) {
		return ErrUnauthorized
	}
	return nil
}

func authenticate(ctx context.Context, obj interface{}, next graphql.Resolver, on *bool) (res interface{}, err error) {
	if auth.GetUser(ctx) == nil {
		// block calling the next resolver
		return nil, fmt.Errorf("access denied")
	}

	// or let it pass through
	return next(ctx)
}

func (r *Resolver) prepareBigQuery(ctx context.Context, bq models.NewBigQuery, group string) (models.BigqueryMetadata, error) {
	if err := r.ensureGroupOwnsGCPProject(ctx, group, bq.ProjectID); err != nil {
		return models.BigqueryMetadata{}, err
	}

	metadata, err := r.bigquery.TableMetadata(ctx, bq.ProjectID, bq.Dataset, bq.Table)
	if err != nil {
		return models.BigqueryMetadata{}, fmt.Errorf("trying to fetch metadata on table %v, but it does not exist in %v.%v",
			bq.Table, bq.ProjectID, bq.Dataset)
	}

	switch metadata.TableType {
	case bigquery.RegularTable:
	case bigquery.ViewTable:
		fallthrough
	case bigquery.MaterializedView:
		if err := r.accessMgr.AddToAuthorizedViews(ctx, bq.ProjectID, bq.Dataset, bq.ProjectID, bq.Dataset, bq.Table); err != nil {
			return models.BigqueryMetadata{}, err
		}
	default:
		return models.BigqueryMetadata{}, fmt.Errorf("unsupported table type: %v", metadata.TableType)
	}

	return metadata, nil
}
