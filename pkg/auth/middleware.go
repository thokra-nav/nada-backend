package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/navikt/nada-backend/pkg/openapi"
	log "github.com/sirupsen/logrus"
)

func MockJWTValidatorMiddleware() openapi.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			teams := []string{
				"team",
			}
			if mockTeam := r.Header.Get("X-Mock-Team"); mockTeam != "" {
				teams[0] = mockTeam
			}
			ctx := context.WithValue(r.Context(), "teams", teams)
			ctx = context.WithValue(ctx, "member_name", "mock_mockerson")
			r = r.WithContext(context.WithValue(ctx, "preferred_username", "mockuser"))
			next.ServeHTTP(w, r)
		})
	}
}

type teamsCache interface {
	Get(uuid string) (string, bool)
}

func JWTValidatorMiddleware(discoveryURL, clientID string, azureGroups *AzureGroups, teamUUIDs teamsCache) openapi.MiddlewareFunc {
	certificates, err := FetchCertificates(discoveryURL)
	if err != nil {
		log.Fatalf("Fetching signing certificates from IDP: %v", err)
	}
	jwtValidator := JWTValidator(certificates, clientID)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(openapi.CookieAuthScopes) == nil {
				next.ServeHTTP(w, r)
				return
			}

			var claims jwt.MapClaims

			token := jwtauth.TokenFromCookie(r)

			_, err := jwt.ParseWithClaims(token, &claims, jwtValidator)
			if err != nil {
				log.Debugf("parsing token: %v", err)
				w.WriteHeader(http.StatusForbidden)
				_, err = fmt.Fprintf(w, "unauthorized access")
				if err != nil {
					log.Errorf("Writing http response: %v", err)
				}
				return
			}

			email := strings.ToLower(claims["preferred_username"].(string))
			exp := int(claims["exp"].(float64))

			ctx := r.Context()
			ctx = context.WithValue(ctx, "preferred_username", email)
			ctx = context.WithValue(ctx, "token_expiry", exp)
			ctx = context.WithValue(ctx, "member_name", "user:"+email)

			groups, err := azureGroups.GetGroupsForUser(ctx, token, email)
			if err != nil {
				log.Errorf("getting groups for user: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err = fmt.Fprintf(w, "Unauthorized access: %s", err.Error())
				if err != nil {
					log.Errorf("Writing http response: %v", err)
				}
				return
			}

			teams := make([]string, 0)

			for _, uuid := range groups {
				if uid, found := teamUUIDs.Get(uuid); found {
					teams = append(teams, uid)
				}
			}

			r = r.WithContext(context.WithValue(ctx, "teams", teams))

			next.ServeHTTP(w, r)
		}
	}
}

func JWTValidator(certificates map[string]CertificateList, audience string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		var certificateList CertificateList
		var kid string
		var ok bool

		if claims, ok := token.Claims.(*jwt.MapClaims); !ok {
			return nil, fmt.Errorf("unable to retrieve claims from token")
		} else {
			if valid := claims.VerifyAudience(audience, true); !valid {
				return nil, fmt.Errorf("the token is not valid for this application")
			}
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if kid, ok = token.Header["kid"].(string); !ok {
			return nil, fmt.Errorf("field 'kid' is of invalid type %T, should be string", token.Header["kid"])
		}

		if certificateList, ok = certificates[kid]; !ok {
			return nil, fmt.Errorf("kid '%s' not found in certificate list", kid)
		}

		for _, certificate := range certificateList {
			return certificate.PublicKey, nil
		}

		return nil, fmt.Errorf("no certificate candidates for kid '%s'", kid)
	}
}