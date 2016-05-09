// +build windows

package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	// "github.com/nanobox-io/nanobox/config"
)

// PrivilegeExec runs a command, but assumes your already running as adminsitrator
func PrivilegeExec(command, msg string) {
	fmt.Printf(stylish.Bullet(msg))

	//
	cmd := exec.Command(os.Args[0], strings.Split(command, " ")...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		lumber.Fatal("[commands/commands_windows]", err.Error())
	}
}
