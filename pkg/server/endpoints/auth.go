package endpoints

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zarinit-routers/router-server/internal/user"
	"github.com/zarinit-routers/router-server/pkg/server/middleware"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			log.Error("Failed to bind JSON", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !user.CheckPassword(loginRequest.Password) {
			log.Error("Password is incorrect")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		token, err := GenerateToken()
		if err != nil {
			log.Error("Failed to generate token", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 10).Unix(),
			"iat": time.Now().Unix(),
		})

	return token.SignedString(middleware.GetSecurityKey())
}
