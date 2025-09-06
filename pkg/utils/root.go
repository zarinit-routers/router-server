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

// Execute basic implementation: runs command without bash -c
// Kept here to satisfy interface, actual extended errors in execute.go
func Execute(name string, args ...string) (string, error) {
	return ExecuteErr(name, args...)
}
