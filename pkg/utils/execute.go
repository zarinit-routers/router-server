package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

// ExecuteErr runs a command via bash -c, logs the command and its output, and returns stdout as string
func ExecuteErr(name string, args ...string) (string, error) {
	cmdStr := strings.Join(append([]string{name}, args...), " ")
	log.Info("Executing command", "cmd", cmdStr)

	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())
	if outStr != "" {
		log.Info("Command stdout", "output", outStr)
	}
	if errStr != "" {
		log.Info("Command stderr", "output", errStr)
	}
	if err != nil {
		return outStr, fmt.Errorf("command failed: %w", err)
	}
	return outStr, nil
}
