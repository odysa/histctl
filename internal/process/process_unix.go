//go:build darwin || linux

package process

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsRunning checks whether a process with the given name is running.
func IsRunning(processName string) (bool, error) {
	cmd := exec.Command("pgrep", "-x", processName)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("check process %q: %w", processName, err)
	}
	return strings.TrimSpace(string(output)) != "", nil
}
