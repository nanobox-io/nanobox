package dev

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processDevStop ...
type processDevStop struct {
	control      processor.ProcessControl
}

//
func init() {
	processor.Register("dev_stop", stopFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func stopFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevStop{control: control}, nil
}

//
func (dev processDevStop) Results() processor.ProcessControl {
	return dev.control
}

//
func (dev processDevStop) Process() error {
	// set the process mode to dev
	// which will allow isolation of containers
	dev.control.Env = "dev"

	return processor.Run("app_stop", dev.control)
}
