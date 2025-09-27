package dhcp

import (
	"fmt"

	"github.com/zarinit-routers/router-server/pkg/kea"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func GetStatus(_ models.JSONMap) (any, error) {
	status, err := kea.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("failed get dhcp server status: %s", err)
	}
	return status, nil
}
func Enable(_ models.JSONMap) (any, error) {
	if err := kea.Enable(); err != nil {
		return nil, fmt.Errorf("failed start dhcp server: %s", err)
	}
	return models.JSONMap{
		"enabled": true,
	}, nil
}

func Disable(_ models.JSONMap) (any, error) {
	if err := kea.Disable(); err != nil {
		return nil, fmt.Errorf("failed stop dhcp server: %s", err)
	}
	return models.JSONMap{
		"enabled": false,
	}, nil
}
