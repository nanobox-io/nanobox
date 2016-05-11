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
	config ProcessConfig
}

func init() {
	Register("dev_console", devConsoleFunc)
}

func devConsoleFunc(conf ProcessConfig) (Processor, error) {
	return devConsole{conf}, nil
}

func (self devConsole) Results() ProcessConfig {
	return self.config
}

func (self devConsole) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	name := self.config.Meta["name"]
	if name == "" {
		name = "build"
	}
	container, err := docker.GetContainer(fmt.Sprintf("%s-%s", util.AppName(), name))
	if err == nil {
		name = container.ID
	}

	command := []string{"exec", "-u", "gonano", "-it", name, "/bin/bash"}

	switch {
	case self.config.Meta["working_dir"] != "":
		cd := fmt.Sprintf("cd %s; exec \"${SHELL:-sh}\"", self.config.Meta["working_dir"])
		command = append(command, "-c", cd)
	case self.config.Meta["command"] != "":
		command = append(command, "-c", self.config.Meta["command"])
	}

	cmd := exec.Command("docker", command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
