// +build !windows

package util

import (
	"fmt"
	"os"
	"os/exec"
)

// IsPrivileged will return true if the current process is running under a
// privileged user, like root
func IsPrivileged() bool {

	// Execute a syscall to return the user id. If the user id is 0 then we're
	// running with root escalation.
	if os.Geteuid() == 0 {
		return true
	}

	return false
}

// PrivilegeExec runs a command as sudo
func PrivilegeExec(command string) error {
	//
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v --internal", command))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	return cmd.Run()
}
