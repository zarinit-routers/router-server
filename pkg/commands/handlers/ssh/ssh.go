package ssh

import (
	"fmt"

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
	err := systemctl.Enable(SSHService)
	if err != nil {
		return nil, fmt.Errorf("failed to enable ssh: %s", err)
	}
	return models.JSONMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}

func Disable(_ models.JSONMap) (any, error) {
	err := systemctl.Disable(SSHService)
	if err != nil {
		return nil, fmt.Errorf("failed to disable ssh: %s", err)
	}
	return models.JSONMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}
