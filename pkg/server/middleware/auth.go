package middleware

import (
	"fmt"
	"net/http"

	l "github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var log *l.Logger

func init() {
	log = l.WithPrefix("Middleware")
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			logBadAuth(fmt.Errorf("header Authorization is empty"))
			return
		}
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (any, error) {
			return GetSecurityKey(), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			logBadAuth(fmt.Errorf("token parse error: %v", err))
			return
		}
		log.Info("Authentication successful", "token", token)
		c.Next()
	}
}

func GetSecurityKey() []byte {
	return []byte(viper.GetString("jwt-security-key"))
}

func logBadAuth(err error) {
	log.Error("Attempt to access without authorization", "error", err, "status", http.StatusUnauthorized)
}
