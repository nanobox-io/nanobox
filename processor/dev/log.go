package dev

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processDevLog ...
type processDevLog struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_log", devLogFn)
}

//
func devLogFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevLog{control}, nil
}

//
func (devLog processDevLog) Results() processor.ProcessControl {
	return devLog.control
}

//
func (devLog processDevLog) Process() error {
	// set the process mode to dev
	devLog.control.Env = "dev"

	// some messaging about the logging??

	return processor.Run("mist_log", devLog.control)
}
