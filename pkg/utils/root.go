package utils

import (
	"fmt"
	"os"
)

var ErrNotRootUser = fmt.Errorf("operation requires root privileges")

func CheckRoot() error {
	if os.Geteuid() != 0 {
		return ErrNotRootUser
	}
	return nil
}
