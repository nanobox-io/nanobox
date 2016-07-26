package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDevDestroy ...
type processDevDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_destroy", simDestroyFn)
}

//
func simDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevDestroy{control}, nil
}

//
func (simDestroy processDevDestroy) Results() processor.ProcessControl {
	return simDestroy.control
}

//
func (simDestroy processDevDestroy) Process() error {
	simDestroy.control.Env = "sim"

	// we need the vm to be up and running
	if err := processor.Run("provider_setup", simDestroy.control); err != nil {
		return err
	}

	// remove all the services (platform/service/code)
	if err := simDestroy.removeServices(); err != nil {
		return err
	}

	// teardown the app
	return processor.Run("env_teardown", simDestroy.control)
}

// removeServices gets all the services in the app and remove them
func (simDestroy processDevDestroy) removeServices() error {
	bucket := fmt.Sprintf("%s_sim", config.AppID())
	services, err := data.Keys(bucket)
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}

	// go through the services and run a service destroy on each
	simDestroy.control.Display(stylish.Bullet("Removing Services"))
	simDestroy.control.DisplayLevel++
	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(config.AppID(), service, &svc)
			simDestroy.control.Meta["name"] = service
			err := processor.Run("service_destroy", simDestroy.control)
			if err != nil {
				simDestroy.control.Display(stylish.Warning("one of the services did not uninstall:\n%s", err.Error()))
				// continue on to the next one.
				// we should continue trying to remove services
			}
		}
	}
	simDestroy.control.DisplayLevel--
	return nil
}
