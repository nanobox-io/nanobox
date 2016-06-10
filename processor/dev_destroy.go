package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type devDestroy struct {
	control ProcessControl
}

func init() {
	Register("dev_destroy", devDestroyFunc)
}

func devDestroyFunc(control ProcessControl) (Processor, error) {
	return devDestroy{control}, nil
}

func (destroy devDestroy) Results() ProcessControl {
	return destroy.control
}

func (destroy devDestroy) Process() error {

	if err := Run("dev_setup", destroy.control); err != nil {
		return err
	}

	// remove all the services (platform/service/code)
	if err := destroy.removeServices(); err != nil {
		return err
	}

	// teardown the app
	if err := Run("app_teardown", destroy.control); err != nil {
		return err
	}

	if err := destroy.removeMounts(); err != nil {
		return err
	}

	// potentially destroy the provider
	if err := destroy.destroyProvider(); err != nil {
		return err
	}

	return nil
}

// get all the services in the app
// and remove them
func (destroy devDestroy) removeServices() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}
	destroy.control.Display(stylish.Bullet("Removing Services"))
	destroy.control.DisplayLevel++
	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(util.AppName(), service, &svc)
			destroy.control.Meta["name"] = service
			err := Run("service_destroy", destroy.control)
			if err != nil {
				destroy.control.Display(stylish.Warning("one of the services did not uninstall:\n%s", err.Error()))
				// continue on to the next one.
				// we should continue trying to remove services
			}
		}
	}
	destroy.control.DisplayLevel--
	return nil
}

// removeMounts will add the shares and mounts for this app
func (destroy devDestroy) removeMounts() error {

	// unmount the engine if it's a local directory
	if util.EngineDir() != "" {
		src := util.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), util.AppName())

		// unmount the share on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			return err
		}

		// remove the share on the workstation
		if err := provider.RemoveShare(src, dst); err != nil {
			return err
		}
	}

	// unmount the app src
	src := util.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), util.AppName())

	// unmount the share on the provider
	if err := provider.RemoveMount(src, dst); err != nil {
		return err
	}

	// remove the share on the workstation
	if err := provider.RemoveShare(src, dst); err != nil {
		return err
	}

	return nil
}

// destroyProvider destroys the provider if there are no remaining apps
func (destroy devDestroy) destroyProvider() error {
	// fetch all of the apps
	keys, err := data.Keys("apps")
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		// if no other apps exist in container
		if err := Run("provider_destroy", destroy.control); err != nil {
			return err
		}
	}
	return nil
}
