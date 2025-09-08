package endpoints

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/router-server/pkg/cloud"
	config "github.com/zarinit-routers/router-server/pkg/cloud/config"
)

func GetConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConnectionConfig()
		c.JSON(http.StatusOK, gin.H{
			"config": conf,
		})
	}
}

type UpdateConfigRequest struct {
	Passphrase     string `json:"passphrase"`
	OrganizationId string `json:"organizationId"`
}

func UpdateConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateConfigRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Failed bind JSON", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		conf := config.GetConnectionConfig()

		if req.Passphrase != "" {
			conf.SetPassphrase(req.Passphrase)
		}
		if req.OrganizationId != "" {
			conf.SetOrganizationId(req.OrganizationId)
		}

		c.JSON(http.StatusOK, gin.H{
			"config": conf,
		})
	}
}
func GetCloudStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := cloud.GetStatus()
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	}
}
