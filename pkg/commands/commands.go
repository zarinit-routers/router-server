package commands

import (
	c "github.com/zarinit-routers/commands"
	"github.com/zarinit-routers/router-server/pkg/commands/handlers/modems"
	"github.com/zarinit-routers/router-server/pkg/commands/handlers/ssh"
	"github.com/zarinit-routers/router-server/pkg/commands/handlers/system"
	"github.com/zarinit-routers/router-server/pkg/commands/handlers/timezone"
	"github.com/zarinit-routers/router-server/pkg/models"
)

type CommandHandler func(models.JSONMap) (any, error)

var implementedCommands = map[string]CommandHandler{
	// timezone
	c.CommandTimezoneGet.String(): timezone.Get,
	c.CommandTimezoneSet.String(): timezone.Set,
	// system
	"v1/system/get-os-info":     system.GetOSInfo,
	"v1/system/get-device-info": system.GetDeviceInfo,
	"v1/system/reboot":          system.Reboot,
	// ssh
	"v1/ssh/enable":     ssh.Enable,
	"v1/ssh/disable":    ssh.Disable,
	"v1/ssh/get-status": ssh.GetStatus,
	// modems
	"v1/modems/list":       modems.List,
	"v1/modems/enable":     modems.Enable,
	"v1/modems/disable":    modems.Disable,
	"v1/modems/get-signal": modems.GetSignal,
}

func CheckCommand(command string) (CommandHandler, error) {
	handler, ok := implementedCommands[command]
	if !ok {
		return nil, NotImplementedErr{Command: command}
	}
	return handler, nil
}
