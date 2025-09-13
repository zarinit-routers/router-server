package diagnostics

import (
	"fmt"

	"github.com/zarinit-routers/router-server/pkg/cli/nslookup"
	"github.com/zarinit-routers/router-server/pkg/cli/ping"
	"github.com/zarinit-routers/router-server/pkg/cli/traceroute"
	"github.com/zarinit-routers/router-server/pkg/models"
)

var (
	ErrAddressNotSpecified = fmt.Errorf("address key not specified")
)

func Ping(args models.JSONMap) (any, error) {

	addr, ok := args["address"].(string)
	if !ok {
		return nil, ErrAddressNotSpecified
	}
	output, err := ping.Run(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to run ping: %s", err)
	}
	return output, nil
}

func Traceroute(args models.JSONMap) (any, error) {

	addr, ok := args["address"].(string)
	if !ok {
		return nil, ErrAddressNotSpecified
	}
	output, err := traceroute.Run(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to run traceroute: %s", err)
	}
	return output, nil
}

func Nslookup(args models.JSONMap) (any, error) {

	addr, ok := args["address"].(string)
	if !ok {
		return nil, ErrAddressNotSpecified
	}
	output, err := nslookup.Run(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to run nslookup: %s", err)
	}
	return output, nil
}
