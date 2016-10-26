// +build windows

package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	// ensure the command is prepared to be used in powershell escalation
	command = preparePrivilegeCmd(command)

	// add the --internal flag if the command is nanobox
	if strings.HasPrefix(command, "nanobox ") {
		command = fmt.Sprintf("%s --internal", command)
	}

	// fetch the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to determine current working directory: %s", cwd)
	}

	// we need to cd into the current directory before running our command
	command = fmt.Sprintf("%s & cd %s & %s", filepath.VolumeName(cwd), cwd, command)

	// generate the powershell process
	process := fmt.Sprintf("& {Start-Process 'cmd' -ArgumentList '/c %s' -Verb RunAs -Wait}", command)

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

// make sure the command is escaped and prepared to be used in powershell
func preparePrivilegeCmd(command string) string {

	// return the command if an .exe wasn't provided
	if !strings.Contains(command, ".exe") || strings.HasSuffix(command, ".exe\""){
		return command
	}

	// split the command into two parts
	parts := strings.Split(command, ".exe")

	r, err := regexp.Compile("([a-zA-Z]+)$")
	if err != nil {
		return command
	}

	// extract the executable from the command
	executable := r.FindString(parts[0])

	// generate a new command without the absolute path to the .exe
	return fmt.Sprintf("%s%s", executable, parts[1])
}
