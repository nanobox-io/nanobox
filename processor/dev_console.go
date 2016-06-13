package processor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/util"
)

// processDevConsole ...
type processDevConsole struct {
	control ProcessControl
}

//
func init() {
	Register("dev_console", devConsoleFunc)
}

//
func devConsoleFunc(conf ProcessControl) (Processor, error) {
	return processDevConsole{conf}, nil
}

//
func (devConsole processDevConsole) Results() ProcessControl {
	return devConsole.control
}

//
func (devConsole processDevConsole) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", devConsole.control)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	name := devConsole.control.Meta["name"]
	if name == "" {
		name = "build"
	}

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", util.AppName(), name))
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
