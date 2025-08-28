package commands

import "fmt"

type NotImplementedErr struct {
	Command string
}

func (e NotImplementedErr) Error() string {
	return fmt.Sprintf("command %q not implemented", e.Command)
}
