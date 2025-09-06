package endpoints

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	mgmt "github.com/zarinit-routers/router-server/pkg/commands/handlers/timezone"
)

type timezoneRequest struct {
	Timezone string `json:"timezone"`
}

type timezoneResponse struct {
	Timezone string `json:"timezone"`
}

func TimezoneHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tz, err := mgmt.GetTimeZone()
		if err != nil {
			log.Errorf("failed to get timezone: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, timezoneResponse{Timezone: tz})
	}
}

func SetTimezoneHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req timezoneRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := mgmt.SetTimeZone(req.Timezone); err != nil {
			log.Errorf("failed to set timezone: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, timezoneResponse{Timezone: req.Timezone})
	}
}
