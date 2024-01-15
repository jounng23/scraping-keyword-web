package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(provideGinEngine),
		fx.Invoke(registerHTTPService, startHTTPServer),
	)
	app.Run()
}

func provideGinEngine() *gin.Engine {
	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))
	return r
}

func registerHTTPService(g *gin.Engine) {
}

func startHTTPServer(lifecycle fx.Lifecycle, g *gin.Engine) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := fmt.Sprintf("%d", "8000")
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
