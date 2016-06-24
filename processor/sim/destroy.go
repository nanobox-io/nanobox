package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/netfs"
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
	
	if err := processor.Run("provider_setup", simDestroy.control); err != nil {
		return err
	}

	// remove all the services (platform/service/code)
	if err := simDestroy.removeServices(); err != nil {
		return err
	}

	// teardown the app
	if err := processor.Run("app_teardown", simDestroy.control); err != nil {
		return err
	}

	if err := simDestroy.removeMounts(); err != nil {
		return err
	}

	// potentially destroy the provider
	if err := simDestroy.destroyProvider(); err != nil {
		return err
	}

	return nil
}

// removeServices gets all the services in the app and remove them
func (simDestroy processDevDestroy) removeServices() error {
	bucket := fmt.Sprintf("%s_%s", config.AppName(), simDestroy.control.Env)
	services, err := data.Keys(bucket)
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}
	simDestroy.control.Display(stylish.Bullet("Removing Services"))
	simDestroy.control.DisplayLevel++
	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(config.AppName(), service, &svc)
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

// removeMounts will add the shares and mounts for this app
func (simDestroy processDevDestroy) removeMounts() error {

	// unmount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

		// unmount the share on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			return err
		}

		// remove the share on the workstation
		if err := simDestroy.removeShare(src, dst); err != nil {
			return err
		}
	}

	// unmount the app src
	src := config.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.AppName())

	// unmount the share on the provider
	if err := provider.RemoveMount(src, dst); err != nil {
		return err
	}

	// remove the share on the workstation
	if err := simDestroy.removeShare(src, dst); err != nil {
		return err
	}

	return nil
}

// destroyProvider destroys the provider if there are no remaining apps
func (simDestroy processDevDestroy) destroyProvider() error {
	// fetch all of the apps
	keys, err := data.Keys("apps")
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		// if no other apps exist in container
		if err := processor.Run("provider_destroy", simDestroy.control); err != nil {
			return err
		}
	}
	return nil
}

// removeShare removes a previously exported share
func (simDestroy processDevDestroy) removeShare(src, dst string) error {

	// we don't really care what mount-type the user has configured, we need
	// to remove any shares

	// first we check netfs
	if netfs.Exists(src) {
		control := processor.ProcessControl{
			Env:      simDestroy.control.Env,
			Verbose:      simDestroy.control.Verbose,
			DisplayLevel: simDestroy.control.DisplayLevel,
			Meta: map[string]string{
				"path": src,
			},
		}

		if err := processor.Run("share_netfs_remove", control); err != nil {
			return err
		}
	}

	// now provider native
	if provider.HasMount(src) {
		if err := provider.RemoveMount(src, dst); err != nil {
			return err
		}
	}

	return nil
}
