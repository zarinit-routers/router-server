package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Cors() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     viper.GetStringSlice("client.addresses"),
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	log.Info("CORS middleware initialized", "AllowOrigins", cfg.AllowOrigins)
	return cors.New(cfg)
}
