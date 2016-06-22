package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
)

// processDevStart ...
type processDevStart struct {
	control      processor.ProcessControl
}

//
func init() {
	processor.Register("dev_start", startFn)
}

// TODO: do some control validation check on the meta for the flags and make sure
// they work
func startFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevStart{control: control}, nil
}

//
func (dev processDevStart) Results() processor.ProcessControl {
	return dev.control
}

//
func (dev processDevStart) Process() error {
	// set the process mode to dev
	// which will allow isolation of containers
	dev.control.Env = "dev"

	// defer the clean up so if we exit early the cleanup will always happen
	defer func() {
		if err := processor.Run("share_teardown", dev.control); err != nil {
			fmt.Println("teardown broke")
			fmt.Println(err)

			// this is bad, really bad...
			// we should probably print a pretty message explaining that the app
			// was left in a bad state and needs to be reset
			return 
		}
	}()

	// get the vm and app up.
	if err := processor.Run("share_setup", dev.control); err != nil {
		return err
	}

	// startDataServices will start all data services
	if err := processor.Run("service_start_all", dev.control); err != nil {
		return err
	}

	if err := dev.watchMist(); err != nil {
		return err
	}

	return nil
}

func (dev *processDevStart) watchMist() error {
	// output some message
	dev.control.Display("while this console is open your dev env will be available")
	dev.control.Display("attempting to connect to live streaming logs")
	dev.control.Display("Next: run a build 'nanobox build'")
	dev.control.Display("Then: open a console and start coding 'nanobox dev console'")
	
	// tail the mist logs
	return processor.Run("mist_log", dev.control)
}
