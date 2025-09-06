package timezone

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/router-server/pkg/utils"
)

func Get(_ map[string]any) (map[string]any, error) {
	info, err := getInfo()
	if err != nil {
		log.Errorf("failed to get timedatectl info: %v", err)
		return nil, err
	}
	return map[string]any{"timezone": info.GetTimeZone()}, nil
}

func Set(params map[string]any) (map[string]any, error) {
	if err := utils.CheckRoot(); err != nil {
		log.Warn("Operation is not allowed for non-root users")
		return nil, err
	}

	tz, _ := params["timezone"].(string)
	if tz == "" {
		return nil, fmt.Errorf("timezone is required")
	}

	_, err := utils.Execute("timedatectl", "set-timezone", tz)
	if err != nil {
		log.Errorf("failed to set timezone: %v", err)
		return nil, err
	}
	return map[string]any{"timezone": tz}, nil
}
