package processor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util"
)

type devConsole struct {
	control ProcessControl
}

func init() {
	Register("dev_console", devConsoleFunc)
}

func devConsoleFunc(conf ProcessControl) (Processor, error) {
	return devConsole{conf}, nil
}

func (self devConsole) Results() ProcessControl {
	return self.control
}

func (self devConsole) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.control)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	name := self.control.Meta["name"]
	if name == "" {
		name = "build"
	}

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", util.AppName(), name))
	if err == nil {
		name = container.ID
	}

	shell := self.control.Meta["shell"]
	if shell == "" {
		shell = "bash"
	}

	command := []string{"exec", "-u", "gonano", "-it", name, "/bin/bash"}

	switch {
	case self.control.Meta["working_dir"] != "":
		cd := fmt.Sprintf("cd %s; exec \"%s\"", self.control.Meta["working_dir"], shell)
		command = append(command, "-c", cd)
	case self.control.Meta["command"] != "":
		command = append(command, "-c", self.control.Meta["command"])
	}

	cmd := exec.Command("docker", command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
