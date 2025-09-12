package sim

import (
	"fmt"

	"github.com/zarinit-routers/router-server/pkg/cli/mmcli"
	"github.com/zarinit-routers/router-server/pkg/models"
)

func Get(args models.JSONMap) (any, error) {
	name, ok := args["sim"].(string)
	if !ok {
		return nil, fmt.Errorf("field 'sim' is required and must be a string")
	}
	sim, err := mmcli.GetSim(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get sim %q: %s", name, err)
	}
	return models.JSONMap{"sim": sim}, nil
}
