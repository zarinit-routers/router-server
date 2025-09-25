package hostapd

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli/systemctl"
	"github.com/zarinit-routers/router-server/pkg/cli/hostapd/config"
)

type Type string

const (
	WIFI2_4G Type = "2.4G"
	WIFI5G   Type = "5G"
)

const (
	Service2_4G = "z-hostapd-2"
	Service5G   = "z-hostapd-5"

	ConfigFile2_4G = "/etc/hostapd/z-hostapd-2.conf"
	ConfigFile5G   = "/etc/hostapd/z-hostapd-5.conf"
)

func (t Type) Enable() error {
	var service string
	switch t {
	case WIFI2_4G:
		service = Service2_4G
	case WIFI5G:
		service = Service5G
	default:
		return fmt.Errorf("unknown wifi type %q", t)
	}

	err := systemctl.Enable(service)
	if err != nil {
		return fmt.Errorf("failed enable service: %s", err)
	}
	return nil
}

func (t Type) Disable() error {
	var service string
	switch t {
	case WIFI2_4G:
		service = Service2_4G
	case WIFI5G:
		service = Service5G
	default:
		return fmt.Errorf("unknown wifi type %q", t)
	}

	err := systemctl.Disable(service)
	if err != nil {
		return fmt.Errorf("failed disable service: %s", err)
	}
	return nil
}
func (t Type) IsActive() bool {
	var service string
	switch t {
	case WIFI2_4G:
		service = Service2_4G
	case WIFI5G:
		service = Service5G
	default:
		log.Error("No such wifi type, can't return error in IsActive function", "type", t)
		return false
	}

	return systemctl.IsActive(service)
}

func (t Type) GetConfig() (*config.Config, error) {
	var confFile string
	switch t {
	case WIFI2_4G:
		confFile = ConfigFile2_4G
	case WIFI5G:
		confFile = ConfigFile5G
	default:
		return nil, fmt.Errorf("unknown wifi type %q", t)
	}
	conf, err := config.FromFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("failed acquire configuration: %s", err)
	}
	applyHWMode(conf, t)
	return conf, nil
}

func (t Type) ChangeSSID(ssid string) error {
	return t.updateConfig(func(c *config.Config) error {
		return c.SetSSID(ssid)
	})
}
func (t Type) ChangePassword(pass string) error {
	return t.updateConfig(func(c *config.Config) error {
		return c.SetPassphrase(pass)
	})
}
func (t Type) ChangeVisibility(visible bool) error {
	return t.updateConfig(func(c *config.Config) error {
		c.SetSSIDVisibility(visible)
		return nil
	})
}
func (t Type) ChangeChannel(channel int) error {
	return t.updateConfig(func(c *config.Config) error {
		return c.SetChannel(channel)
	})
}

func (t Type) updateConfig(confFunc func(*config.Config) error) error {
	conf, err := t.GetConfig()
	if err != nil {
		return fmt.Errorf("failed acquire configuration: %s", err)
	}
	if err := confFunc(conf); err != nil {
		return fmt.Errorf("failed change configuration: %s", err)
	}
	if err := conf.Write(); err != nil {
		return fmt.Errorf("failed set config")
	}
	return nil
}

func applyHWMode(c *config.Config, t Type) {
	switch t {
	case WIFI2_4G:
		if c.GetHWMode() != config.HWMode2_4G {
			c.SetHWMode(config.HWMode2_4G)
		}
	case WIFI5G:
		if c.GetHWMode() != config.HWMode5G {
			c.SetHWMode(config.HWMode5G)
		}
	}
}
