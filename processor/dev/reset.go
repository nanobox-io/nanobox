package dev

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDevReset ...
type processDevReset struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_reset", devResetFn)
}

//
func devResetFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevReset{control: control}, nil
}

//
func (devReset processDevReset) Results() processor.ProcessControl {
	return devReset.control
}

//
func (devReset processDevReset) Process() error {

	if err := devReset.resetCounters(); err != nil {
		return err
	}

	return nil
}

// resetCounters resets all the counters associated with all apps
func (devReset processDevReset) resetCounters() error {

	apps, err := data.Keys("apps")

	if err != nil {
		return err
	}

	for _, app := range apps {
		// reset the app dev usage counter
		counter.Reset(app + "_dev")
	}

	return nil
}
