package dockerexec

import (
  "io"
  "bytes"
  "errors"

  "github.com/nanobox-io/golang-docker-client"
)

type Cmd struct {
  Id      string
  Path    string
  Payload string
  Out     bytes.Buffer
  Stdout  io.Writer
}

// Run builds a command and executes within the context of a docker container
func (cmd *Cmd) Run() error {
  // assemble the full command to run within the hooks dir
  run := []string{"/opt/nanobox/hooks/" + cmd.Path, cmd.Payload}

  // start the exec
	exec, hj, err := docker.ExecStart(cmd.Id, run, false, true, true)
	if err != nil {
		return err
	}

  // stream the output
	if err := docker.ExecPipe(hj, nil, &cmd.Out, cmd.Stdout); err != nil {
    return err
  }

  // let's see if we can inspect a bit
	data, err := docker.ExecInspect(exec.ID)
	if err != nil {
		return err
	}

  // was the exit code bad?
	if data.ExitCode != 0 {
		return errors.New("bad exit code")
	}

	return nil
}

// Output returns the output from the command
func (cmd *Cmd) Output() string {
  return cmd.Out.String()
}

// Command generates a new Cmd struct
func Command(id string, path string, payload string) *Cmd {
  return &Cmd{
    Id: id,
    Path: path,
    Payload: payload,
  }
}
