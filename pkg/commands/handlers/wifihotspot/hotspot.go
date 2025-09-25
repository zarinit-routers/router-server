package wifihotspot

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/zarinit-routers/cli/iw"
	"github.com/zarinit-routers/router-server/pkg/cli/hostapd"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func Enable(_ models.JSONMap) (any, error) {
	err := hostapd.WIFI2_4G.Enable()
	if err != nil {
		return nil, fmt.Errorf("failed enable wifi hotspot: %s", err)
	}

	return models.JSONMap{"enabled": true}, nil
}

func Disable(_ models.JSONMap) (any, error) {
	err := hostapd.WIFI2_4G.Disable()
	if err != nil {
		return nil, fmt.Errorf("failed disable wifi hotspot: %s", err)
	}

	return models.JSONMap{"enabled": false}, nil
}

func GetStatus(_ models.JSONMap) (any, error) {
	conf, err := hostapd.WIFI2_4G.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed get config: %s", err)
	}

	return models.JSONMap{
		"enabled":  hostapd.WIFI2_4G.IsActive(),
		"ssid":     conf.GetSSID(),
		"password": conf.GetPassphrase(),
		"channel":  conf.GetChannel(),
		"hidden":   conf.GetSSIDVisibility(),
	}, nil
}

func SetSSID(args models.JSONMap) (any, error) {
	ssid, ok := args["ssid"].(string)
	if !ok {
		return nil, fmt.Errorf("ssid not specified")
	}
	err := hostapd.WIFI2_4G.ChangeSSID(ssid)
	if err != nil {
		return nil, fmt.Errorf("failed change ssid: %s", err)
	}

	return models.JSONMap{"ssid": ssid}, nil
}

func SetSSIDVisibility(args models.JSONMap) (any, error) {
	hidden, ok := args["hidden"].(bool)
	if !ok {
		return nil, fmt.Errorf("visibility not specified (key 'hidden' of type bool)")
	}
	err := hostapd.WIFI2_4G.ChangeVisibility(!hidden)
	if err != nil {
		return nil, fmt.Errorf("failed change visibility: %s", err)
	}
	return models.JSONMap{"hidden": hidden}, nil
}

func SetPassword(args models.JSONMap) (any, error) {
	password, ok := args["password"].(string)
	if !ok {
		return nil, fmt.Errorf("password not specified")
	}

	err := hostapd.WIFI2_4G.ChangePassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed set password: %s", err)
	}
	return models.JSONMap{"password": password}, nil
}

func SetChannel(args models.JSONMap) (any, error) {
	channel, ok := args["channel"].(int)
	if !ok {
		return nil, fmt.Errorf("channel not specified")
	}

	err := hostapd.WIFI2_4G.ChangeChannel(channel)
	if err != nil {
		return nil, fmt.Errorf("failed set channel: %s", err)
	}
	return models.JSONMap{"channel": channel}, nil
}
func GetConnectedClients(_ models.JSONMap) (any, error) {
	clients, err := iw.GetConnectedDevices(viper.GetString("wifi-hotspot.interface"))
	if err != nil {
		return nil, fmt.Errorf("failed get connected clients: %s", err)
	}
	return models.JSONMap{"clients": clients}, nil
}
