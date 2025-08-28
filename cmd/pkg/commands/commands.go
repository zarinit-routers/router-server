package commands

import (
	"github.com/zarinit-routers/router-server/cmd/pkg/commands/handlers/timezone"
)

type JsonMap = map[string]any

type CommandHandler func(JsonMap) (JsonMap, error)

var commands = map[string]CommandHandler{
	"v1/timezone/get": timezone.Get,
	"v1/timezone/set": timezone.Set,
}

func CheckCommand(command string) (CommandHandler, error) {
	handler, ok := commands[command]
	if !ok {
		return nil, NotImplementedErr{Command: command}
	}
	return handler, nil
}
