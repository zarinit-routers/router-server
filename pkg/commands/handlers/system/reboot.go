package system

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func Reboot(_ models.JSONMap) (any, error) {
	go func() {
		time.Sleep(5 * time.Second)
		err := cli.ExecuteErr("reboot")
		if err != nil {
			log.Error("Failed reboot", "error", err)
		}
	}()

	return models.JSONMap{
		"message": "Rebooting...",
	}, nil

}
