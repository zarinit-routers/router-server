package server

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/pkg/server/endpoints"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	cfg := cors.Config{
		AllowOrigins:     viper.GetStringSlice("client.addresses"),
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(cfg))

	api := r.Group("/api")
	{
		api.POST("/cmd", endpoints.CommandHandler())
	}

	log.Info("HTTP server initialized", "AllowOrigins", cfg.AllowOrigins)
	return r
}
