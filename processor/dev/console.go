package dev

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
)

// processDevConsole ...
type processDevConsole struct {
	control   processor.ProcessControl
	name      string
	container string
	command   string
	cwd       string
	shell     string
}

func init() {
	processor.Register("dev_console", devConsoleFn)
}

//
func devConsoleFn(control processor.ProcessControl) (processor.Processor, error) {
	devConsole := &processDevConsole{control: control}
	return devConsole, devConsole.validateMeta()
}

func (devConsole processDevConsole) Results() processor.ProcessControl {
	return devConsole.control
}

//
func (devConsole processDevConsole) Process() error {

	// setup the environment (boot vm)
	if err := processor.Run("provider_setup", devConsole.control); err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	//
	if container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", config.AppName(), devConsole.container)); err == nil {
		devConsole.container = container.ID
	}

	// this is the default command to run in the container
	command := []string{"exec", "-u", "gonano", "-it", devConsole.container, "/bin/bash"}

	// check to see if there are any optional meta arguments that need to be handled
	switch {

	// if a current working directory (cwd) is provided then modify the command to
	// change into that directory before executing
	case devConsole.cwd != "":
		command = append(command, "-c", fmt.Sprintf("cd %s; exec \"%s\"", devConsole.cwd, devConsole.shell))

	// if a command is provided then modify the command to exec that command after
	// running the base command
	case devConsole.command != "":
		command = append(command, "-c", devConsole.command)
	}

	//
	cmd := exec.Command("docker", command...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//
	return cmd.Run()
}

// validateMeta validates that the required metadata exists
func (devConsole *processDevConsole) validateMeta() error {

	// set optional meta values
	devConsole.command = devConsole.control.Meta["command"]
	devConsole.cwd = devConsole.control.Meta["cwd"]

	// set container; if no container is provided default to "build"
	devConsole.container = devConsole.control.Meta["container"]
	if devConsole.container == "" {
		devConsole.container = "build"
	}

	// set shell; if no shell is provided default to "bash"
	devConsole.shell = devConsole.control.Meta["shell"]
	if devConsole.shell == "" {
		devConsole.shell = "bash"
	}

	return nil
}
