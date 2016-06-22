package share

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
)

// processShareConsole ...
type processShareConsole struct {
	control   processor.ProcessControl
	name      string
	container string
	command   string
	cwd       string
	shell     string
}

func init() {
	processor.Register("share_console", shareConsoleFn)
}

//
func shareConsoleFn(control processor.ProcessControl) (processor.Processor, error) {
	shareConsole := &processShareConsole{control: control}
	return shareConsole, shareConsole.validateMeta()
}

func (shareConsole processShareConsole) Results() processor.ProcessControl {
	return shareConsole.control
}

//
func (shareConsole processShareConsole) Process() error {

	// setup the environment (boot vm)
	if err := processor.Run("provider_setup", shareConsole.control); err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	//
	id := fmt.Sprintf("nanobox-%s-%s-%s", config.AppName(), shareConsole.control.Env, shareConsole.container)
	if container, err := docker.GetContainer(id); err == nil {
		shareConsole.container = container.ID
	}

	// this is the default command to run in the container
	command := []string{"exec", "-u", "gonano", "-it", shareConsole.container, "/bin/bash"}

	// check to see if there are any optional meta arguments that need to be handled
	switch {

	// if a current working directory (cwd) is provided then modify the command to
	// change into that directory before executing
	case shareConsole.cwd != "":
		command = append(command, "-c", fmt.Sprintf("cd %s; exec \"%s\"", shareConsole.cwd, shareConsole.shell))

	// if a command is provided then modify the command to exec that command after
	// running the base command
	case shareConsole.command != "":
		command = append(command, "-c", shareConsole.command)
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
func (shareConsole *processShareConsole) validateMeta() error {

	// set optional meta values
	shareConsole.command = shareConsole.control.Meta["command"]
	shareConsole.cwd = shareConsole.control.Meta["cwd"]

	// set container; if no container is provided default to "build"
	shareConsole.container = shareConsole.control.Meta["container"]
	if shareConsole.container == "" {
		shareConsole.container = fmt.Sprintf("nanobox-%s-build", config.AppName())
	}

	// set shell; if no shell is provided default to "bash"
	shareConsole.shell = shareConsole.control.Meta["shell"]
	if shareConsole.shell == "" {
		shareConsole.shell = "bash"
	}

	return nil
}
