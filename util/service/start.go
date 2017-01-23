package service

import (
	"fmt"
	"os/exec"
)


func Start(name string) error {
	cmd := startCmd(name)
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		fmt.Errorf("out: %s, err: %s", out, err)
	}
	return nil
}
