// +build windows

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nanobox-io/nanobox-golang-stylish"
)

// privilegeExec runs a command but assumes 
// your already running as adminsitrator
func privilegeExec(command, msg string) {
	fmt.Printf(stylish.Bullet(msg))

	//
	cmd := exec.Command(os.Args[0], strings.Split(command, " ")...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		Config.Fatal("[commands/halt]", err.Error())
	}
}
