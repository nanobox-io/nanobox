package env

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox/models"
	// "github.com/nanobox-io/nanobox/processors/provider"
)

// Console ...
type Console struct {
	Component models.Component
	Command   string
	Cwd       string
	Shell     string
}

//
func (console Console) Run() error {
	// set the default shell
	if console.Shell == "" {
		console.Shell = "bash"
	}

	// // setup the environment (boot vm)
	// providerSetup := provider.Setup{}
	// if err := providerSetup.Run(); err != nil {
	// 	return err
	// }

	// this is the default command to run in the container
	cmd := []string{
		"docker",
		"exec",
		"-u",
		"gonano",
		"-it",
		console.Component.ID,
		"/bin/bash",
	}

	// check to see if there are any optional meta arguments that need to be handled
	switch {

	// if a current working directory (cwd) is provided then modify the command to
	// change into that directory before executing
	case console.Cwd != "":
		cmd = append(cmd, "-c", fmt.Sprintf("cd %s; exec \"%s\"", console.Cwd, console.Shell))

	// if a command is provided then modify the command to exec that command after
	// running the base command
	case console.Command != "":
		cmd = append(cmd, "-c", console.Command)
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	if err := process.Run(); err != nil && err.Error() != "exit status 137" {
		return err
	}

	return nil
}
