package api

import (
	"net/http"

	graphProm "github.com/99designs/gqlgen-contrib/prometheus"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/database"
	"github.com/navikt/nada-backend/pkg/graph"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func New(
	repo *database.Repo,
	gcp graph.GCP,
	oauth2 OAuth2,
	gcpProjects *auth.TeamProjectsUpdater,
	accessMgr graph.AccessManager,
	schemaUpdater graph.SchemaUpdater,
	authMW auth.MiddlewareHandler,
	log *logrus.Logger,
) *chi.Mux {
	corsMW := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	httpAPI := new(oauth2, log.WithField("subsystem", "api"))

	gqlServer := graph.New(repo, gcp, gcpProjects, accessMgr, schemaUpdater, log.WithField("subsystem", "graph"))

	router := chi.NewRouter()
	router.Use(corsMW)
	router.Route("/api", func(r chi.Router) {
		r.Handle("/", playground.Handler("GraphQL playground", "/api/query"))
		r.Handle("/query", authMW(gqlServer))
		r.HandleFunc("/login", httpAPI.Login)
		r.HandleFunc("/oauth2/callback", httpAPI.Callback)
		r.HandleFunc("/logout", httpAPI.Logout)
	})
	router.Route("/internal", func(r chi.Router) {
		r.Handle("/metrics", prometheusGenerator(repo))
	})

	return router
}

func prometheusGenerator(repo *database.Repo) http.Handler {
	registry := prometheus.NewRegistry()
	graphProm.RegisterOn(registry)
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{ReportErrors: true}))

	if err := repo.Register(registry); err != nil {
		panic(err)
	}
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
