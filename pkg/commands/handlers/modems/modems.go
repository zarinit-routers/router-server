package modems

import (
	"fmt"

	"github.com/zarinit-routers/router-server/pkg/cli/mmcli"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func ListModems(_ models.JSONMap) (any, error) {
	modems, err := mmcli.List()
	if err != nil {
		return nil, err
	}
	return models.JSONMap{"modems": modems}, nil
}

var ErrInvalidModemArg = fmt.Errorf("field 'modem' is required and must be a string value")

func EnableModem(args models.JSONMap) (any, error) {
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

func DisableModem(args models.JSONMap) (any, error) {
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
