package dev

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processDevStart ...
type processDevStart struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_start", startFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func startFn(control processor.ProcessControl) (processor.Processor, error) {
	return &processDevStart{control: control}, nil
}

//
func (dev *processDevStart) Results() processor.ProcessControl {
	return dev.control
}

//
func (dev *processDevStart) Process() error {
	// set the process mode to dev
	dev.control.Env = "dev"

	if err := processor.Run("app_start", dev.control); err != nil {
		return err
	}

	// messaging about what happened and next steps

	return nil
}
