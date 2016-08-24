package hookit

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
)

// Cmd ...
type Cmd struct {
	ID      string
	Path    string
	Payload string
	Stdout  io.Writer
	Stderr  io.Writer
}

// Run builds a command and executes within the context of a docker container
func (cmd *Cmd) Run() error {
	lumber.Debug("hookit:Cmd:Run: %s, %s, %s", cmd.ID, cmd.Path, cmd.Payload)

	// assemble the full command to run within the hooks dir
	run := []string{"/opt/nanobox/hooks/" + cmd.Path, cmd.Payload}

	// start the exec
	exec, hj, err := docker.ExecStart(cmd.ID, run, false, true, true)
	if err != nil {
		return err
	}

	// if no streams are given then set a reasonable alternative
	// this will be later used to make the error messaging
	// better
	var buff bytes.Buffer
	if cmd.Stderr == nil {
		cmd.Stderr = &buff
	}

	// stream the output
	if err := docker.ExecPipe(hj, nil, cmd.Stdout, cmd.Stderr); err != nil {
		return err
	}

	// let's see if we can inspect a bit
	data, err := docker.ExecInspect(exec.ID)
	if err != nil {
		return err
	}

	// was the exit code bad?
	if data.ExitCode != 0 {
		// if so use the buffer that may have been assigned to the
		// streams to give message better error handling
		return fmt.Errorf("bad exit code(%d): %s", data.ExitCode, buff.String())
	}

	return nil
}

// Output returns the output from the command
// this mirrors the Output command of the os/exec package
func (cmd *Cmd) Output() (string, error) {
	if cmd.Stdout != nil {
		return "", fmt.Errorf("stdout is already set")
	}

	var buffer bytes.Buffer
	cmd.Stdout = &buffer
  
	err := cmd.Run()
	if err != nil {
		lumber.Error("hookit:Cmd:Run: %s, %s, %s", cmd.ID, cmd.Path, cmd.Payload)
		lumber.Error("hookit:cmd:Output: %s, err: %s", buffer.String(), err.Error())
		err = fmt.Errorf("failed to run hook (%s) in container: %s", cmd.Path, err.Error())
	}

	lumber.Debug("hookit:Cmd:Output: %s", buffer.String())

	return buffer.String(), err
}

// DockerCommand generates a new Cmd struct
func DockerCommand(id string, path string, payload string) *Cmd {
	return &Cmd{
		ID:      id,
		Path:    path,
		Payload: payload,
	}
}
