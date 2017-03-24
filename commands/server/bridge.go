package server

import (
	util_bridge "github.com/nanobox-io/nanobox/util/provider/bridge"
)


type Bridge struct {}

var bridge = &Bridge{
}

func (br *Bridge) Start(config string, resp *Response) error {
	resp = &Response{
		Output: "",
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
		Output: "",
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

	return run("Bridge.Start", config, resp)	
}

func StopBridge() error {
	resp := &Response{}

	return run("Bridge.Stop", "", resp)
}