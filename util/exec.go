package util

import (
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
	var b bytes.Buffer
	err = docker.ExecPipe(hj, nil, &b, &b)
	if err != nil {
		return b.String(), err
	}
	lumber.Debug(" - result from exec:\n%s", b.String())
	data, err := docker.ExecInspect(exec.ID)
	if err != nil {
		return b.String(), err
	}
	if data.ExitCode != 0 {
		return b.String(), badExit
	}
	return b.String(), nil
}
