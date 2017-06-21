package bridge

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/util/config"
)

func Start(conf string) error {
	resp := &Response{}

	return server.ClientRun("Bridge.Start", conf, resp)
}

func (br *Bridge) Start(conf string, resp *Response) error {

	if runningBridge != nil {
		// if asked to start but we are already running
		// lets teardown the old and recreate with the new conf
		br.Stop(conf, resp)
		// return nil
	}

	runningBridge = exec.Command(config.VpnPath(), "--config", conf)
	runningBridge.Stdout = os.Stdout
	runningBridge.Stderr = os.Stderr

	err := runningBridge.Start()
	if err != nil {
		runningBridge = nil
		err = fmt.Errorf("failed to start cmd(%s --config %s): %s", config.VpnPath(), conf, err)
	}
	return err
}
