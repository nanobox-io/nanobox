package env

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Console ...
func Console(componentModel *models.Component, consoleConfig ConsoleConfig) error {
	// set the default shell
	if consoleConfig.Shell == "" {
		consoleConfig.Shell = "bash"
	}

	// setup docker client
	if err := provider.Init(); err != nil {
		return err
	}

	// this is the default command to run in the container
	cmd := []string{
		"docker",
		"exec",
		"-u",
		"gonano",
		"-it",
		componentModel.ID,
		"/bin/bash",
	}

	// check to see if there are any optional meta arguments that need to be handled
	switch {

	// if a current working directory (cwd) is provided then modify the command to
	// change into that directory before executing
	case consoleConfig.Cwd != "":
		cmd = append(cmd, "-c", fmt.Sprintf("cd %s; exec \"%s\"", consoleConfig.Cwd, consoleConfig.Shell))

	// if a command is provided then modify the command to exec that command after
	// running the base command
	case consoleConfig.Command != "":
		cmd = append(cmd, "-c", consoleConfig.Command)
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	if err := process.Run(); err != nil && err.Error() != "exit status 137" && err.Error() != "exit status 130" {
		return err
	}

	return nil
}
