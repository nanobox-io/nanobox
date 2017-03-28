package server

import (
	util_bridge "github.com/nanobox-io/nanobox/util/provider/bridge"
)

type Bridge struct{}

func init() {
	Register(&Bridge{})
}

func (br *Bridge) Start(config string, resp *Response) error {
	resp = &Response{
		Output:   "",
		ExitCode: 0,
	}

	err := util_bridge.Start(config)
	if err != nil {
		resp.ExitCode = 1
	}
	return err
}

func (br *Bridge) Stop(config string, resp *Response) error {
	resp = &Response{
		Output:   "",
		ExitCode: 0,
	}

	err := util_bridge.Stop()
	if err != nil {
		resp.ExitCode = 1
	}
	return err
}

func StartBridge(config string) error {
	resp := &Response{}

	return ClientRun("Bridge.Start", config, resp)
}

func StopBridge() error {
	resp := &Response{}

	return ClientRun("Bridge.Stop", "", resp)
}
