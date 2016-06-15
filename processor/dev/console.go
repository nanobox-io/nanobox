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
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_console", devConsoleFn)
}

//
func devConsoleFn(conf processor.ProcessControl) (processor.Processor, error) {
	return processDevConsole{conf}, nil
}

//
func (devConsole processDevConsole) Results() processor.ProcessControl {
	return devConsole.control
}

//
func (devConsole processDevConsole) Process() error {
	// setup the environment (boot vm)
	err := processor.Run("provider_setup", devConsole.control)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	name := devConsole.control.Meta["name"]
	if name == "" {
		name = "build"
	}

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", config.AppName(), name))
	if err == nil {
		name = container.ID
	}

	shell := devConsole.control.Meta["shell"]
	if shell == "" {
		shell = "bash"
	}

	command := []string{"exec", "-u", "gonano", "-it", name, "/bin/bash"}

	switch {

	//
	case devConsole.control.Meta["working_dir"] != "":
		cd := fmt.Sprintf("cd %s; exec \"%s\"", devConsole.control.Meta["working_dir"], shell)
		command = append(command, "-c", cd)

		//
	case devConsole.control.Meta["command"] != "":
		command = append(command, "-c", devConsole.control.Meta["command"])
	}

	cmd := exec.Command("docker", command...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
