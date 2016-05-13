package util

import (
	"os"
	"bytes"
	"errors"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
)

var badExit = errors.New("bad exit code")

func Exec(id, name, payload string) (string, error) {
	lumber.Debug("Execing %s in container %s with a payload of %s", name, id, payload)
	exec, hj, err := docker.ExecStart(id, []string{"/opt/nanobox/hooks/" + name, payload}, false, true, true)
	if err != nil {
		return "", err
	}
	var stdout bytes.Buffer
	err = docker.ExecPipe(hj, nil, &stdout, os.Stdout)
	if err != nil {
		return stdout.String(), err
	}
	lumber.Debug(" - result from exec out:\n%s", stdout.String())
	data, err := docker.ExecInspect(exec.ID)
	if err != nil {
		return stdout.String(), err
	}
	if data.ExitCode != 0 {
		return stdout.String(), badExit
	}
	return stdout.String(), nil
}
