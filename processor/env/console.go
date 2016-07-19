package env

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
)

// processEnvConsole ...
type processEnvConsole struct {
	control   processor.ProcessControl
	name      string
	container string
	command   string
	cwd       string
	shell     string
}

func init() {
	processor.Register("env_console", envConsoleFn)
}

//
func envConsoleFn(control processor.ProcessControl) (processor.Processor, error) {
	envConsole := &processEnvConsole{control: control}
	return envConsole, envConsole.validateMeta()
}

func (envConsole processEnvConsole) Results() processor.ProcessControl {
	return envConsole.control
}

//
func (envConsole processEnvConsole) Process() error {

	// setup the environment (boot vm)
	if err := processor.Run("provider_setup", envConsole.control); err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		return err
	}

	//
	id := fmt.Sprintf("nanobox_%s_%s_%s", config.AppName(), envConsole.control.Env, envConsole.container)
	if container, err := docker.GetContainer(id); err == nil {
		envConsole.container = container.ID
	}

	// this is the default command to run in the container
	command := []string{"exec", "-u", "gonano", "-it", envConsole.container, "/bin/bash"}

	// check to see if there are any optional meta arguments that need to be handled
	switch {

	// if a current working directory (cwd) is provided then modify the command to
	// change into that directory before executing
	case envConsole.cwd != "":
		command = append(command, "-c", fmt.Sprintf("cd %s; exec \"%s\"", envConsole.cwd, envConsole.shell))

	// if a command is provided then modify the command to exec that command after
	// running the base command
	case envConsole.command != "":
		command = append(command, "-c", envConsole.command)
	}

	//
	cmd := exec.Command("docker", command...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//
	
	if err := cmd.Run(); err != nil && err.Error() != "exit status 137" {
		return err
	}
	return nil
}

// validateMeta validates that the required metadata exists
func (envConsole *processEnvConsole) validateMeta() error {

	// set optional meta values
	envConsole.command = envConsole.control.Meta["command"]
	envConsole.cwd = envConsole.control.Meta["cwd"]

	// set container; if no container is provided default to "build"
	envConsole.container = envConsole.control.Meta["container"]
	if envConsole.container == "" {
		envConsole.container = fmt.Sprintf("nanobox_%s_build", config.AppName())
	}

	// set shell; if no shell is provided default to "bash"
	envConsole.shell = envConsole.control.Meta["shell"]
	if envConsole.shell == "" {
		envConsole.shell = "bash"
	}

	return nil
}
