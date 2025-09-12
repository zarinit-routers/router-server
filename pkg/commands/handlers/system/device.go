package system

import (
	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func GetDeviceInfo(_ models.JSONMap) (any, error) {
	return models.JSONMap{
		"manufacturer":    viper.GetString("device.manufacturer"),
		"model":           viper.GetString("device.model"),
		"modelVersion":    viper.GetString("device.model-version"),
		"firmwareVersion": viper.GetString("device.firmware-version"),
		"id":              viper.GetString("device.id"),
	}, nil
}
