// +build windows

package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IsPrivileged will return true if the current process is running as the
// Administrator
func IsPrivileged() bool {
	// Running "net session" will return "Access is denied." if the terminal
	// process was not run as Administrator
	cmd := exec.Command("net", "session")
	output, err := cmd.CombinedOutput()

	// if there was an error, we'll short-circuit and return false
	if err != nil {
		return false
	}

	// return false if we find Access is denied in the output
	if bytes.Contains(output, []byte("Access is denied.")) {
		return false
	}

	// if the previous checks didn't fail, then we must be the Administrator
	return true
}

// PrivilegeExec will run the requested command in a powershell as the Administrative user
func PrivilegeExec(command string) error {

	// Windows is tricky. Unfortunately we can't just prefix the command with sudo
	// Instead, we have to use powershell to create a profile, and then create
	// a process within powershell requesting Administrative permissions.
	//
	// Generating the command is complicated.
	// The following resources were used as documentation for the logic below:
	// https://msdn.microsoft.com/en-us/powershell/scripting/core-powershell/console/powershell.exe-command-line-help
	// http://ss64.com/ps/start-process.html
	// http://www.howtogeek.com/204088/how-to-use-a-batch-file-to-make-powershell-scripts-easier-to-run/

	// The process is constructed by passing the executable as a single argument
	// and the argument list as a space-delimited string in a single argument.
	//
	// Since the command is provided as a space-delimited string containing both
	// the executable and the argument list (just like a command would be entered
	// on the command prompt), we need to pop off the executable.

	executable, arguments := splitExecutableAndArgs(command)

	// generate the powershell process
	process := fmt.Sprintf("& {Start-Process '%s' -ArgumentList '%s --internal' -Verb RunAs -Wait}", executable, arguments)

	// now we can generate a command to exec
	cmd := exec.Command("PowerShell.exe", "-NoProfile", "-Command", process)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// PowerShell will run a specified command in a powershell and return the result
func PowerShell(command string) ([]byte, error) {

	process := fmt.Sprintf("& {%s}", command)

	cmd := exec.Command("PowerShell.exe", "-NoProfile", "-Command", process)

	return cmd.CombinedOutput()
}

// extracts the executable from the args
func splitExecutableAndArgs(cmd string) (executable, args string) {

	if strings.Contains(cmd, ".exe") {
		// split the command by the .exe extension
		parts := strings.Split(cmd, ".exe ")
		// the first item is the executable
		executable = fmt.Sprintf("%s.exe", parts[0])
		// the second item are the args
		args = parts[1]
		
		return
	}
	
	// split the command by spaces
	parts := strings.Split(cmd, " ")
	// extract the executable (the first item)
	executable = parts[0]
	// the remaining are the args
	args = strings.Join(parts[1:], " ")

	return
}
