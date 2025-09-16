package endpoints

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zarinit-routers/router-server/pkg/models"
)

type LoginRequest struct {
	Name     string `json:"username" binding:"required"`
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

		user, err := models.GetUserByUsername(loginRequest.Name)
		if err != nil {
			log.Error("Failed to get user", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		err = user.CheckPassword(loginRequest.Password)
		if err != nil {
			log.Error("Failed to check password", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, err := GenerateToken(user)
		if err != nil {
			log.Error("Failed to generate token", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 10).Unix(),
			"iat": time.Now().Unix(),
			"sub": user.ID,
		})

	return token.SignedString([]byte("jwt-security-key"))
}
