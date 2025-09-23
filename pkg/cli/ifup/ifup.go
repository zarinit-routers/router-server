package ifup

import (
	"github.com/zarinit-routers/cli"
)

func Up(ifname string) error {
	return cli.ExecuteErr("ifup", ifname)
}
