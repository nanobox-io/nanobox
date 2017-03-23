package server

import (
	util_bridge "github.com/nanobox-io/nanobox/util/provider/bridge"
)


type Bridge struct {}

var bridge = &Bridge{
}

func (br *Bridge) Start(config string, resp *Response) error {
	return util_bridge.Start(config)
}

func (br *Bridge) Stop(config string, resp *Response) error {
	return util_bridge.Stop()
}

func StartBridge(config string) error {
	resp := &Response{}

	return run("Bridge.Start", config, resp)	
}

func StopBridge() error {
	resp := &Response{}

	return run("Bridge.Stop", "", resp)
}