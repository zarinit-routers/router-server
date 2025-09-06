package timezone

import (
	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/router-server/pkg/utils"
)

func Get() (string, error) {
	info, err := getInfo()
	if err != nil {
		log.Errorf("failed to get timedatectl info: %v", err)
		return "", err
	}
	return info.GetTimeZone(), nil
}

func Set(tz string) error {
	if err := utils.CheckRoot(); err != nil {
		log.Warn("Operation is not allowed for non-root users")
		return err
	}
	_, err := utils.Execute("timedatectl", "set-timezone", tz)
	if err != nil {
		log.Errorf("failed to set timezone: %v", err)
		return err
	}
	return nil
}
