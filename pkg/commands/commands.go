package commands

import (
	c "github.com/zarinit-routers/commands"
	"github.com/zarinit-routers/router-server/pkg/commands/handlers/timezone"
)

type JsonMap = map[string]any

type CommandHandler func(JsonMap) (JsonMap, error)

var implementedCommands = map[string]CommandHandler{
	c.CommandTimezoneGet.String(): timezone.Get,
	c.CommandTimezoneSet.String(): timezone.Set,
}

func CheckCommand(command string) (CommandHandler, error) {
	handler, ok := implementedCommands[command]
	if !ok {
		return nil, NotImplementedErr{Command: command}
	}
	return handler, nil
}
