package sim

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processSimStop ...
type processSimStop struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_stop", stopFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func stopFn(control processor.ProcessControl) (processor.Processor, error) {
	return processSimStop{control: control}, nil
}

//
func (sim processSimStop) Results() processor.ProcessControl {
	return sim.control
}

//
func (sim processSimStop) Process() error {
	// set the process mode to sim
	// which will allow isolation of containers
	sim.control.Env = "sim"

	return processor.Run("app_stop", sim.control)
}
