package util

import (
	"errors"
	"bytes"
	"github.com/nanobox-io/golang-docker-client"
)

var badExit = errors.New("bad exit code")

func Exec(id, name, payload string) (string, error) {
	exec, hj, err := docker.ExecStart(id, []string{"/opt/bin/"+name, payload}, false, true, true)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = docker.ExecPipe(hj, nil, &b, &b)
	if err != nil {
		return b.String(), err
	}
	data, err := docker.ExecInspect(exec.ID)	
	if err != nil {
		return b.String(), err
	}
	if data.ExitCode != 0 {
		return b.String(), badExit
	}
	return b.String(), nil
}