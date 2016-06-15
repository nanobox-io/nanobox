// +build !windows

package util

import (
	"fmt"
	"os"
	"os/exec"
)

// PrivilegeExec runs a command as sudo
func PrivilegeExec(command string) error {
	//
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v", command))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
