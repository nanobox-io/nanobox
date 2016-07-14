package sim

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processSimStart ...
type processSimStart struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_start", startFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func startFn(control processor.ProcessControl) (processor.Processor, error) {
	return processSimStart{control: control}, nil
}

//
func (sim processSimStart) Results() processor.ProcessControl {
	return sim.control
}

//
func (sim processSimStart) Process() error {
	// set the process mode to sim
	sim.control.Env = "sim"

	return processor.Run("app_start", sim.control)
}
