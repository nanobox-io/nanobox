package bridge

import (
	"fmt"
	"os/exec"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/util/config"
)
type stream struct {

}

func (s stream) Write(p []byte) (n int, err error) {
	lumber.Info("bridge: %s", p)
	return len(p), nil
}

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
	runningBridge.Stdout = stream{}
	runningBridge.Stderr = stream{}

	err := runningBridge.Start()
	if err != nil {
		runningBridge = nil
		err = fmt.Errorf("failed to start cmd(%s --config %s): %s", config.VpnPath(), conf, err)
	}
	return err
}
