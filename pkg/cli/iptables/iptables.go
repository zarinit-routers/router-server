package iptables

import (
	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
	"github.com/zarinit-routers/router-server/pkg/utils"
)

// Returns port forwarding logs. Only if user is root
func PortForwarding() (string, error) {
	if err := utils.CheckRoot(); err != nil {
		log.Warnf("Can't get port forwarding logs: %s", err)
		return "", err
	}
	output, err := cli.Execute("iptables", "-t", "nat", "-L", "PREROUTING", "-n", "-v")
	return string(output), err
}
