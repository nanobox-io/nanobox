package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
)

// processSimStart ...
type processSimStart struct {
	control      processor.ProcessControl
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
	// which will allow isolation of containers
	sim.control.Env = "sim"

	// defer the clean up so if we exit early the cleanup will always happen
	defer func() {
		if err := processor.Run("env_teardown", sim.control); err != nil {
			fmt.Println("teardown broke")
			fmt.Println(err)

			return 
		}
	}()

	// get the vm and app up.
	if err := processor.Run("env_setup", sim.control); err != nil {
		return err
	}

	// startDataServices will start all data services
	if err := processor.Run("service_start_all", sim.control); err != nil {
		return err
	}

	return sim.watchMist()
}

func (sim *processSimStart) watchMist() error {
	// output some message
	sim.control.Display("             _  _ ____ _  _ ____ ___  ____ _  _")
	sim.control.Display(`             |\ | |__| |\ | |  | |__) |  |  \/`)
	sim.control.Display(`             | \| |  | | \| |__| |__) |__| _/\_`)
	sim.control.Display("")
	sim.control.Display("----------------------------------------------------------------")
	sim.control.Display("while this console is open your sim env will be available")
	sim.control.Display("attempting to connect to live streaming logs")
	sim.control.Display("Next: run a build 'nanobox build'")
	sim.control.Display("Then: open a console and start coding 'nanobox dev console'")
	sim.control.Display("----------------------------------------------------------------")
	
	// tail the mist logs
	return processor.Run("mist_log", sim.control)
}
