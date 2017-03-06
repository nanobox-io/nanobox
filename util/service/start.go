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

	if !running(name) {
		return fmt.Errorf("service start was successful but the service is not running")
	}
	return nil
}
