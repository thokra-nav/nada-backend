package auth

import (
	"context"
	"net/http"
)

var MockUser = User{
	Name:  "Anderson, Mock",
	Email: "mock.anderson@email.com",
	GoogleGroups: Groups{
		{
			Name:  "team",
			Email: "team@nav.no",
		},
		{
			Name:  "nada",
			Email: "nada@nav.no",
		},
		{
			Name:  "aura",
			Email: "aura@nav.no",
		},
	},
}

var MockProjectIDs = []string{"team-dev-1337", "team-prod-4242"}

var MockTeamProjectsUpdater = TeamProjectsUpdater{
	teamProjects: map[string][]string{
		"team@nav.no": MockProjectIDs,
		"nada@nav.no": {"dataplattform-dev-9da3"},
		"aura@nav.no": {"aura-dev-d9f5"},
	},
}

type MiddlewareHandler func(http.Handler) http.Handler

func MockJWTValidatorMiddleware() MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-NO-AUTH") != "" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, &MockUser)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
