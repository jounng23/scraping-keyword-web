package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"scraping-keyword-web/backend/internal/di/configfx"
	"scraping-keyword-web/backend/internal/di/crawlerfx"
	"scraping-keyword-web/backend/internal/di/dbfx"
	"scraping-keyword-web/backend/internal/di/keywordfx"
	"scraping-keyword-web/backend/internal/di/routerfx"
	"scraping-keyword-web/backend/internal/di/userfx"
	"scraping-keyword-web/backend/internal/di/userkeywordfx"
	"scraping-keyword-web/backend/internal/server/router"
	"scraping-keyword-web/backend/pkg/config"
)

func main() {
	app := fx.New(
		configfx.Initialize(".env"),
		dbfx.Module,
		crawlerfx.Module,
		userfx.Module,
		userkeywordfx.Module,
		keywordfx.Module,
		routerfx.Module,
		fx.Provide(provideGinEngine),
		fx.Invoke(registerHTTPService,
			startHTTPServer),
	)
	app.Run()
}

func provideGinEngine() *gin.Engine {
	r := gin.New()
	r.Use(CORS())
	return r
}

func registerHTTPService(g *gin.Engine,
	router router.Router) {
	api := g.Group("/api/v1")
	router.Register(api)
}

func startHTTPServer(lifecycle fx.Lifecycle, g *gin.Engine) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := fmt.Sprintf("%d", config.ServerConfig().APIPort)
				log.Info().Msgf("listen HTTP on port: %s", port)
				go func() {
					server := http.Server{
						Addr:    ":" + port,
						Handler: g,
					}
					if err := server.ListenAndServe(); err != nil {
						log.Error().Msgf("failed to listen and serve from server: %v", err.Error())
					}
				}()
				return nil
			},
			OnStop: func(context.Context) error {
				log.Info().Msg("service stopped")
				return nil
			},
		},
	)
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("%s:%v", config.ServerConfig().Host, config.ServerConfig().WebAppPort))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
