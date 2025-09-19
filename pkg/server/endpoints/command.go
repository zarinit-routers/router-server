package endpoints

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/router-server/pkg/commands"
)

type CommandRequest struct {
	Command string         `json:"command" binding:"required"`
	Args    map[string]any `json:"args"`
}

type CommandResponse struct {
	Data         any    `json:"data"`
	CommandError string `json:"commandError"`
	RequestError string `json:"requestError"`
}

func requestError(err error) CommandResponse {
	return CommandResponse{
		RequestError: err.Error(),
	}
}
func commandError(err error) CommandResponse {
	return CommandResponse{
		CommandError: err.Error(),
	}
}

func CommandHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req CommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Failed bind JSON", "error", err)
			c.JSON(http.StatusBadRequest, requestError(err))
			return
		}

		cmd, err := commands.CheckCommand(req.Command)
		if err != nil {
			log.Error("Failed check command", "error", err)
			c.JSON(http.StatusBadRequest, requestError(err))
			return
		}

		data, err := cmd(req.Args)
		if err != nil {
			log.Error("Failed execute command", "error", err)
			c.JSON(http.StatusBadRequest, commandError(err))
			return
		}

		c.JSON(http.StatusOK, CommandResponse{
			Data: data,
		})

	}
}
