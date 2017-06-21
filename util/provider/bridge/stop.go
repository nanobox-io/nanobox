package bridge

import (
	"github.com/nanobox-io/nanobox/commands/server"
)

func Stop() error {
	resp := &Response{}

	return server.ClientRun("Bridge.Stop", "", resp)
}

func (br *Bridge) Stop(config string, resp *Response) error {
	if runningBridge == nil {
		return nil
	}

	if err := runningBridge.Process.Kill(); err != nil {
		return err
	}

	if err := runningBridge.Wait(); err != nil {
		// it gets a signal but it shows up as an error
		// we dont want that
		return nil
	}

	// if we killed it and released the resources
	// remove running bridge
	runningBridge = nil
	return nil
}
