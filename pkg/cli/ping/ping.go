package ping

import (
	"fmt"
	"strconv"

	"github.com/zarinit-routers/cli"
)

const (
	DefaultCount = 5
)

func Run(addr string) (string, error) {
	if addr == "" {
		return "", fmt.Errorf("address is empty")
	}

	output, err := cli.Execute(
		"ping",
		addr,
		"-c",
		strconv.Itoa(DefaultCount),
	)
	return string(output), err
}
