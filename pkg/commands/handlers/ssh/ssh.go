package ssh

import (
	"github.com/zarinit-routers/cli/systemctl"
	"github.com/zarinit-routers/router-server/pkg/models"
)

const SSHService = "sshd"

func GetStatus(_ models.JsonMap) (any, error) {
	return models.JsonMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}
func Enable(_ models.JsonMap) (any, error) {
	systemctl.Enable(SSHService)
	return models.JsonMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}

func Disable(_ models.JsonMap) (any, error) {
	systemctl.Disable(SSHService)
	return models.JsonMap{
		"enabled": systemctl.IsActive(SSHService),
	}, nil
}
