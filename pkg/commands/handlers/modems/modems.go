package modems

import (
	"fmt"

	"github.com/zarinit-routers/router-server/pkg/cli/mmcli"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func List(_ models.JSONMap) (any, error) {
	modems, err := mmcli.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list modems: %s", err)
	}
	return models.JSONMap{"modems": modems}, nil
}

var ErrInvalidModemArg = fmt.Errorf("field 'modem' is required and must be a string value")

func Enable(args models.JSONMap) (any, error) {
	name, ok := args["modem"].(string)
	if !ok {
		return nil, ErrInvalidModemArg
	}
	modem, err := mmcli.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed get modem %q: %s", name, err)
	}
	err = modem.Enable()
	if err != nil {
		return nil, fmt.Errorf("failed enable modem %q: %s", name, err)
	}
	return models.JSONMap{"modem": name}, nil
}

func Disable(args models.JSONMap) (any, error) {
	name, ok := args["modem"].(string)
	if !ok {
		return nil, ErrInvalidModemArg
	}
	modem, err := mmcli.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed get modem %q: %s", name, err)
	}
	err = modem.Disable()
	if err != nil {
		return nil, fmt.Errorf("failed disable modem %q: %s", name, err)
	}
	return models.JSONMap{"modem": name}, nil
}

func GetSignal(args models.JSONMap) (any, error) {
	name, ok := args["modem"].(string)
	if !ok {
		return nil, ErrInvalidModemArg
	}
	modem, err := mmcli.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed get modem %q: %s", name, err)
	}
	signal, err := modem.GetSignal()
	if err != nil {
		return nil, fmt.Errorf("failed get signal for modem %q: %s", name, err)
	}
	return models.JSONMap{"signal": signal}, nil
}
