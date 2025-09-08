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

	authMiddleware := func(c *gin.Context) {
		log.Warn("Auth middleware not implemented")
		c.Next()
	}

	api := r.Group("/api")
	{
		api.POST("/cmd", authMiddleware, endpoints.CommandHandler())
	}
	cloud := api.Group("/cloud")
	{
		cloud.Use(authMiddleware)
		cloud.GET("/config", endpoints.GetConfigHandler())
		cloud.POST("/config", endpoints.UpdateConfigHandler())
		cloud.GET("/status", endpoints.GetCloudStatusHandler())
	}

	log.Info("HTTP server initialized", "AllowOrigins", cfg.AllowOrigins)
	return r
}
