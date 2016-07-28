package env

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDestroy ...
type processDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("env_destroy", destroyFn)
}

//
func destroyFn(control processor.ProcessControl) (processor.Processor, error) {
	envDestroy := &processDestroy{control}
	return envDestroy, envDestroy.validateMeta()
}


func (destroy *processDestroy) validateMeta() error {
	if destroy.control.Env == "" {
		return fmt.Errorf("Env not set")
	}

	if destroy.control.Meta["app_name"] == "" {
		destroy.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppID(), destroy.control.Env)
	}

	return nil
}

//
func (destroy processDestroy) Results() processor.ProcessControl {
	return destroy.control
}

//
func (destroy *processDestroy) Process() error {

	// we need the vm to be up and running
	if err := processor.Run("provider_setup", destroy.control); err != nil {
		return err
	}

	// remove the dev container if there is one
	// but dont catch any errors because there
	// may not be a container
	destroy.removeDev()

	// remove all the services (platform/service/code)
	if err := destroy.removeServices(); err != nil {
		return err
	}

	// remove the app and its meta data
	if err := destroy.destroyApp(); err != nil {
		return err
	}

	// remove all dns entries for this app
	if err := processor.Run("env_dns_remove_all", destroy.control); err != nil {
		return err		
	}

	// destroy the mounts
	return destroy.destroyMounts()
}

// remove the development container if one exists
// if not dont complain
func (destroy *processDestroy) removeDev() {
	name := fmt.Sprintf("nanobox_%s", destroy.control.Meta["app_name"])

	docker.ContainerRemove(name)
}

// removeServices gets all the services in the app and remove them
func (destroy processDestroy) removeServices() error {

	services, err := data.Keys(destroy.control.Meta["app_name"])
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}

	destroy.control.Display(stylish.Bullet("Removing Services"))
	destroy.control.DisplayLevel++

	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(config.AppID(), service, &svc)
			destroy.control.Meta["name"] = service
			err := processor.Run("service_destroy", destroy.control)
			if err != nil {
				destroy.control.Display(stylish.Warning("one of the services did not uninstall(%s):\n%s", service, err.Error()))
				// continue on to the next one.
				// we should continue trying to remove services
			}
		}
	}

	destroy.control.DisplayLevel--

	return nil
}


// destroyApp tears down the app when it's not being used
func (destroy *processDestroy) destroyApp() error {

	// establish a local app lock to ensure we're the only ones bringing down the
	// app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// the app package expects a name. not an app_name
	destroy.control.Meta["name"] = destroy.control.Meta["app_name"] 

	// stop all data services
	if err := processor.Run("app_destroy", destroy.control); err != nil {
		return err
	}

	return nil
}

// destroyMounts removes the environments mounts from the provider
func (destroy *processDestroy) destroyMounts() error {
	// we will not be tearing down the mounts currently
	// this is because they are required for production builds
	return nil
	// return processor.Run("app_unmount", destroy.control)
}
