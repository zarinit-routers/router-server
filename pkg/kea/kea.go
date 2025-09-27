package kea

import (
	"github.com/zarinit-routers/cli/systemctl"
)

const (
	Service = "kea-dhcp4.service"
)

type Status struct {
	Enabled bool `json:"enabled"`
	// Config  any
}

func GetStatus() (*Status, error) {

	return &Status{
		Enabled: systemctl.IsActive(Service),
	}, nil
}

func Enable() error {
	return systemctl.Enable(Service)
}
func Disable() error {
	return systemctl.Disable(Service)
}
