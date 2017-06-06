package api

import (
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
)

const JWKS_URI = "https://mineflip.auth0.com/.well-known/jwks.json"
const AUTH0_API_ISSUER = "https://mineflip.auth0.com/"

var AUTH0_API_AUDIENCE = []string{"32BbtHrb7MbfBgTHoDGmGuv6loMkTtvA"}

func parseJWT(r *http.Request) map[string]interface{} {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: JWKS_URI})
	audience := AUTH0_API_AUDIENCE

	configuration := auth0.NewConfiguration(client, audience, AUTH0_API_ISSUER, jose.RS256)
	validator := auth0.NewValidator(configuration)

	token, err := validator.ValidateRequest(r)
	if err != nil {
		return nil
	}

	claims := map[string]interface{}{}
	err = validator.Claims(r, token, &claims)
	if err != nil {
		return nil
	}
	return claims
}
