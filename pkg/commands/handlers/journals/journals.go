package journals

import (
	"fmt"

	"github.com/zarinit-routers/router-server/pkg/cli/iptables"
	"github.com/zarinit-routers/router-server/pkg/cli/journalctl"
	"github.com/zarinit-routers/router-server/pkg/cli/ss"
	"github.com/zarinit-routers/router-server/pkg/models"
)

const (
	JournalCore           = "core"
	JournalSystem         = "system"
	JournalPortForwarding = "port-forwarding"
	JournalConnections    = "connections"
)

func Get(args models.JSONMap) (any, error) {

	journal, ok := args["journal"].(string)
	if !ok {
		return nil, fmt.Errorf("journal key is required and must be a string")
	}

	journalDump := ""
	var err error
	switch journal {
	case JournalCore:
		journalDump, err = journalctl.Core()
		if err != nil {
			return nil, fmt.Errorf("failed to get core journal: %s", err)
		}
	case JournalSystem:
		journalDump, err = journalctl.System()
		if err != nil {
			return nil, fmt.Errorf("failed to get system journal: %s", err)
		}
	case JournalPortForwarding:
		journalDump, err = iptables.PortForwarding()
		if err != nil {
			return nil, fmt.Errorf("failed to get port forwarding journal: %s", err)
		}
	case JournalConnections:
		journalDump, err = ss.Connections()
		if err != nil {
			return nil, fmt.Errorf("failed to get connections journal: %s", err)
		}
	default:
		return nil, fmt.Errorf("invalid journal type")
	}
	return models.JSONMap{"journal": journalDump}, nil
}
