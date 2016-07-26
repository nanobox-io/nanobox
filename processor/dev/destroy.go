package dev

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
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
	processor.Register("dev_destroy", devDestroyFn)
}

//
func devDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevDestroy{control}, nil
}

//
func (devDestroy processDevDestroy) Results() processor.ProcessControl {
	return devDestroy.control
}

//
func (devDestroy processDevDestroy) Process() error {
	devDestroy.control.Env = "dev"

	if err := processor.Run("provider_setup", devDestroy.control); err != nil {
		return err
	}

	// remove the dev container if there is one
	// but dont catch any errors because there
	// may not be a container
	devDestroy.removeDev()

	// remove all the services (platform/service/code)
	if err := devDestroy.removeServices(); err != nil {
		return err
	}

	// teardown the environment and app
	return processor.Run("env_teardown", devDestroy.control)
}

// removeServices gets all the services in the app and remove them
func (devDestroy processDevDestroy) removeServices() error {

	bucket := fmt.Sprintf("%s_dev", config.AppID())
	services, err := data.Keys(bucket)
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}

	devDestroy.control.Display(stylish.Bullet("Removing Services"))
	devDestroy.control.DisplayLevel++

	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(config.AppID(), service, &svc)
			devDestroy.control.Meta["name"] = service
			err := processor.Run("service_destroy", devDestroy.control)
			if err != nil {
				devDestroy.control.Display(stylish.Warning("one of the services did not uninstall:\n%s", err.Error()))
				// continue on to the next one.
				// we should continue trying to remove services
			}
		}
	}

	devDestroy.control.DisplayLevel--

	return nil
}

// remove the development container if one exists
// if not dont complain
func (devDestroy processDevDestroy) removeDev() {
	name := fmt.Sprintf("nanobox_%s_dev", config.AppID())

	docker.ContainerRemove(name)
}
