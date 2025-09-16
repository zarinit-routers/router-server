package server

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/router-server/pkg/server/endpoints"
	"github.com/zarinit-routers/router-server/pkg/server/middleware"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.Cors())

	authMiddleware := middleware.Auth()

	{
		api := r.Group("/api")
		api.POST("/cmd", authMiddleware, endpoints.CommandHandler())

		cloud := api.Group("/cloud")
		{
			cloud.Use(authMiddleware)
			cloud.GET("/config", endpoints.GetConfigHandler())
			cloud.POST("/config", endpoints.UpdateConfigHandler())
			cloud.GET("/status", endpoints.GetCloudStatusHandler())
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", endpoints.LoginHandler())
		}
	}

	log.Info("HTTP server initialized")
	return r
}
