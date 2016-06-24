package dev

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

	// remove all the services (platform/service/code)
	if err := devDestroy.removeServices(); err != nil {
		return err
	}

	// teardown the app
	if err := processor.Run("app_teardown", devDestroy.control); err != nil {
		return err
	}

	if err := devDestroy.removeMounts(); err != nil {
		return err
	}

	// potentially destroy the provider
	if err := devDestroy.destroyProvider(); err != nil {
		return err
	}

	return nil
}

// removeServices gets all the services in the app and remove them
func (devDestroy processDevDestroy) removeServices() error {
	bucket := fmt.Sprintf("%s_%s", config.AppName(), devDestroy.control.Env)
	services, err := data.Keys(bucket)
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}
	devDestroy.control.Display(stylish.Bullet("Removing Services"))
	devDestroy.control.DisplayLevel++
	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(config.AppName(), service, &svc)
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

// removeMounts will add the shares and mounts for this app
func (devDestroy processDevDestroy) removeMounts() error {

	// unmount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

		// unmount the share on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			return err
		}

		// remove the share on the workstation
		if err := devDestroy.removeShare(src, dst); err != nil {
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
	if err := devDestroy.removeShare(src, dst); err != nil {
		return err
	}

	return nil
}

// destroyProvider destroys the provider if there are no remaining apps
func (devDestroy processDevDestroy) destroyProvider() error {
	// fetch all of the apps
	keys, err := data.Keys("apps")
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		// if no other apps exist in container
		if err := processor.Run("provider_destroy", devDestroy.control); err != nil {
			return err
		}
	}
	return nil
}

// removeShare removes a previously exported share
func (devDestroy processDevDestroy) removeShare(src, dst string) error {

	// we don't really care what mount-type the user has configured, we need
	// to remove any shares

	// first we check netfs
	if netfs.Exists(src) {
		control := processor.ProcessControl{
			Env:      devDestroy.control.Env,
			Verbose:      devDestroy.control.Verbose,
			DisplayLevel: devDestroy.control.DisplayLevel,
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
