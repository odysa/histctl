//go:build windows

package process

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsRunning checks whether a process with the given name is running.
func IsRunning(processName string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName), "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("check process %q: %w", processName, err)
	}
	return strings.Contains(strings.ToLower(string(output)), strings.ToLower(processName)), nil
}
