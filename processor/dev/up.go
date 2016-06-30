package dev

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processDevUp ...
type processDevUp struct {
	control      processor.ProcessControl
}

//
func init() {
	processor.Register("dev_up", upFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func upFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevUp{control: control}, nil
}

//
func (dev processDevUp) Results() processor.ProcessControl {
	return dev.control
}

//
func (dev processDevUp) Process() error {
	// set the process mode to dev
	// which will allow isolation of containers
	dev.control.Env = "dev"

	// run a nanobox start
	if err := processor.Run("start", dev.control); err != nil {
		return err
	}

	// run a nanobox build
	if err := processor.Run("build", dev.control); err != nil {
		return err
	}

	// run a dev start
	if err := processor.Run("dev_start", dev.control); err != nil {
		return err
	}

	// run a dev deploy
	if err := processor.Run("dev_deploy", dev.control); err != nil {
		return err
	}

	return nil
}
