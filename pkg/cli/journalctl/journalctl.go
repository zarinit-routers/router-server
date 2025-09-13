package journalctl

import (
	"fmt"

	"github.com/zarinit-routers/cli"
)

const JournalctlExecutable = "journalctl"
const KernelMessagesFlag = "--dmesg" // "-k"
const NoPagerFlag = "--no-pager"

var DefaultLines = 100

func linesFlag(lines int) string {
	return fmt.Sprintf("--lines=%d", lines)
}
func Core() (string, error) {
	output, err := cli.Execute(JournalctlExecutable, KernelMessagesFlag, linesFlag(DefaultLines), NoPagerFlag)
	return string(output), err
}

func System() (string, error) {
	output, err := cli.Execute(JournalctlExecutable, linesFlag(DefaultLines), NoPagerFlag)
	return string(output), err
}
