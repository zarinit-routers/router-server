package ss

import "github.com/zarinit-routers/cli"

func Connections() (string, error) {
	output, err := cli.ExecuteWrap("ss", "-tuln")
	return string(output), err
}
