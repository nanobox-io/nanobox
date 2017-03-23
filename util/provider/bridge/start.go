package bridge

import (
	"fmt"
	"os/exec"
	"os"

	"github.com/nanobox-io/nanobox/util/config"
)

var	runningBridge *exec.Cmd

func Start(conf string) error {
	if runningBridge != nil {
		return nil
	}

	runningBridge = exec.Command(config.VpnPath(), "--config", conf) 
	runningBridge.Stdout = os.Stdout
	runningBridge.Stderr = os.Stderr

	err := runningBridge.Run()
	if err != nil {
		runningBridge = nil
		err = fmt.Errorf("failed to start cmd(%s --config %s): %s",  config.VpnPath(), conf, err)
	}
	return err
}
