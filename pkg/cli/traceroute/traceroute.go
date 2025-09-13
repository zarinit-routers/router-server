package traceroute

import (
	"fmt"

	"github.com/zarinit-routers/cli"
)

func Run(addr string) (string, error) {
	if addr == "" {
		return "", fmt.Errorf("address is empty")
	}
	output, err := cli.Execute("traceroute", addr)
	return string(output), err
}
