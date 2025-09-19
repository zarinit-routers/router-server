package ss

import "github.com/zarinit-routers/cli"

func Connections() (string, error) {
	output, err := cli.Execute("ss", "-tuln")
	return string(output), err
}
