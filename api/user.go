package api

import (
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	jose "gopkg.in/square/go-jose.v2"
)

// User represents authenticated player
type User struct {
	ID   string
	Name string
}

// UserInfoMiddleware binds user object to the request context,
// if it is present in the Authorization header, and valid token.
func UserInfoMiddleware() gin.HandlerFunc {
	conf := struct {
		jwksURI  string
		issuer   string
		audience []string
	}{
		"https://mineflip.auth0.com/.well-known/jwks.json",
		"https://mineflip.auth0.com/",
		[]string{"32BbtHrb7MbfBgTHoDGmGuv6loMkTtvA"},
	}
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: conf.jwksURI})

	configuration := auth0.NewConfiguration(client, conf.audience, conf.issuer, jose.RS256)
	validator := auth0.NewValidator(configuration)
	return func(c *gin.Context) {
		// If no bearer token present, skip
		if c.Request.Header["Authorization"] == nil {
			log.Debug("No auth token, skipping checks")
			c.Next()
			return
		}

		token, err := validator.ValidateRequest(c.Request)
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Error validating auth token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := map[string]interface{}{}
		err = validator.Claims(c.Request, token, &claims)
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Error parsing auth token claims")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := User{ID: claims["sub"].(string), Name: claims["given_name"].(string)}

		c.Set("user", user)

		c.Next()
	}
}
