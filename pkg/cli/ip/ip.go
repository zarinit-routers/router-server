package ip

import (
	"fmt"
	"strings"

	"github.com/zarinit-routers/cli"
)

func GetIP(ifName string) (string, error) {
	output, err := cli.Execute("ip", "--brief", "address", "show", ifName, "primary")
	if err != nil {
		return "", fmt.Errorf("failed execute ip: %v", err)
	}

	words := strings.Fields(string(output))
	if len(words) < 3 {
		return "", fmt.Errorf("failed parse ip output: %v", output)
	}
	return words[2], nil
}
