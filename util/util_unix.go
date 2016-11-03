// +build !windows

package util

import (
	"bufio"
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
	if !sudoExists() {
		fmt.Println("We could not find 'sudo' in your path")
		fmt.Println("Please run the following command in a separate console, then press enter to continue once its complete:")
		fmt.Printf("sudo %v --internal", command)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		return nil
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v --internal", command))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	return cmd.Run()
}

func sudoExists() bool {
	_, err := exec.LookPath("sudo")
	return err == nil
}
