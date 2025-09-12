package ssh

import (
	"github.com/zarinit-routers/cli/systemctl"
	"github.com/zarinit-routers/router-server/pkg/models"
)

const SSHService = "sshd"

func GetStatus(_ models.JSONMap) (any, error) {
	return models.JSONMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}
func Enable(_ models.JSONMap) (any, error) {
	systemctl.Enable(SSHService)
	return models.JSONMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}

func Disable(_ models.JSONMap) (any, error) {
	systemctl.Disable(SSHService)
	return models.JSONMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}
