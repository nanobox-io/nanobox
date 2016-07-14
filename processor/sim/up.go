package sim

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processSimUp ...
type processSimUp struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_up", upFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func upFn(control processor.ProcessControl) (processor.Processor, error) {
	return processSimUp{control: control}, nil
}

//
func (sim processSimUp) Results() processor.ProcessControl {
	return sim.control
}

//
func (sim processSimUp) Process() error {
	// set the process mode to sim
	// which will allow isolation of containers
	sim.control.Env = "sim"

	// run a nanobox start
	if err := processor.Run("start", sim.control); err != nil {
		return err
	}

	// run a nanobox build
	if err := processor.Run("build", sim.control); err != nil {
		return err
	}

	// run a sim start
	if err := processor.Run("sim_start", sim.control); err != nil {
		return err
	}

	// run a sim deploy
	if err := processor.Run("sim_deploy", sim.control); err != nil {
		return err
	}

	return nil
}
