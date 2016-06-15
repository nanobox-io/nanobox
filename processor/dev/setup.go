package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processDevSetup ...
type processDevSetup struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_setup", devSetupFn)
}

//
func devSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	// control.Meta["processDevSetup-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return &processDevSetup{control: control}, nil
}

//
func (devSetup processDevSetup) Results() processor.ProcessControl {
	return devSetup.control
}

//
func (devSetup *processDevSetup) Process() error {

	if err := devSetup.setupProvider(); err != nil {
		return err
	}

	if err := devSetup.setupMounts(); err != nil {
		return err
	}

	if err := devSetup.setupApp(); err != nil {
		return err
	}

	return nil
}

// setupProvider sets up the provider
func (devSetup *processDevSetup) setupProvider() error {

	// let anyone else know we're using the provider
	counter.Increment("provider")

	// establish a global lock to ensure we're the only ones setting up a provider
	// also, we need to ensure the lock is released even if we error
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if err := processor.Run("provider_setup", devSetup.control); err != nil {
		return err
	}

	return nil
}

// setupMounts will add the shares and mounts for this app
func (devSetup *processDevSetup) setupMounts() error {

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

		// first export the share on the workstation
		if err := addShare(src, dst); err != nil {
			return err
		}

		// mount the share on the provider
		if err := addMount(src, dst); err != nil {
			return err
		}
	}

	// mount the app src
	src := config.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.AppName())

	// first export the share on the workstation
	if err := addShare(src, dst); err != nil {
		return err
	}

	// then mount the share on the provider
	if err := addMount(src, dst); err != nil {
		return err
	}

	return nil
}

// setupApp sets up the app plaftorm and data services
func (devSetup *processDevSetup) setupApp() error {

	// let anyone else know we're using the app
	counter.Increment(config.AppName())

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	// setup the app
	if err := processor.Run("app_setup", devSetup.control); err != nil {
		return err
	}

	// clean up after any possible failures in a previous deploy
	if err := processor.Run("service_clean", devSetup.control); err != nil {
		return err
	}

	// setup the platform services
	return processor.Run("platform_setup", devSetup.control)
}

// addShare will add a filesystem share on the host machine
func addShare(src, dst string) error {

	// the mount type is configurable by the user
	mountType := config.Viper().GetString("mount-type")

	// todo: we should display a warning when using native about performance

	// since vm.mount is configurable, it's possible and even likely that a
	// machine may already have mounts configured. For each mount type we'll
	// need to check if an existing mount needs to be undone before continuing
	switch mountType {

	// check to see if netfs is currently configured. If it is then tear it down
	// and build the native share
	case "native":
		if netfs.Exists(src) {
			// netfs was used prior. So we need to tear it down.

			control := ProcessControl{
				Meta: map[string]string{
					"path": src,
				},
			}

			if err := Run("dev_netfs_remove", control); err != nil {
				return err
			}
		}

		// now we let the provider add it's native share
		if err := provider.AddShare(src, dst); err != nil {
			return err
		}

	// check to see if native shares are currently exported. If so,
	// tear down the native share and build the netfs share
	case "netfs":
		if provider.HasShare(src) {
			// native was used prior. So we need to tear it down
			if err := provider.RemoveShare(src, dst); err != nil {
				return err
			}
		}

		control := ProcessControl{
			Meta: map[string]string{
				"path": src,
			},
		}

		if err := Run("dev_netfs_add", control); err != nil {
			return err
		}
	}

	return nil
}

// addMount will mount a share in the nanobox guest context
func addMount(src, dst, string) error {

		// the mount type is configurable by the user
		mountType := config.Viper().GetString("vm.mount")

		switch mountType {

		// build the native mount
		case "native":
			if err := provider.AddMount(src, dst); err != nil {
				return err
			}

		// build the netfs mount
		case "netfs":
			if err := netfs.Mount(src, dst); err != nil {
				return err
			}
		}

		return nil
}
