package nslookup

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
)

func Run(addr string) (string, error) {
	if addr == "" {
		log.Error("Address string is empty")
		return "", fmt.Errorf("address is empty")
	}
	output, err := cli.Execute("nslookup", addr)
	return string(output), err
}
