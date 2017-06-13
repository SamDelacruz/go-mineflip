package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/samdelacruz/go-mineflip/api"
	"github.com/samdelacruz/go-mineflip/hub"
	"gopkg.in/gin-contrib/cors.v1"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT env variable must be set")
	}

	go hub.Run()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	r.Use(cors.New(config))

	r.Use(api.UserInfoMiddleware())

	r.GET("/ws", func(c *gin.Context) {
		hub.HandleWebsocket(c.Writer, c.Request)
	})

	r.POST("/games", api.CreateGameHandler)

	r.GET("/games/:id", api.GetGameHandler)

	r.GET("/games/:id/tiles/:x/:y", api.MoveHandler)

	r.Run()
}
