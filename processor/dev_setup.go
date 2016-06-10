package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
)

type devSetup struct {
	control ProcessControl
}

func init() {
	Register("dev_setup", devSetupFunc)
}

func devSetupFunc(control ProcessControl) (Processor, error) {
	// control.Meta["devSetup-control"]
	// do some control validation
	// check on the meta for the flags and make sure they work

	return &devSetup{control: control}, nil
}

func (setup devSetup) Results() ProcessControl {
	return setup.control
}

func (setup *devSetup) Process() error {

	if err := setup.setupProvider(); err != nil {
		return err
	}

	if err := setup.setupMounts(); err != nil {
		return err
	}

	if err := setup.setupApp(); err != nil {
		return err
	}

	return nil
}

// setupProvider sets up the provider
func (setup *devSetup) setupProvider() error {

	// let anyone else know we're using the provider
	counter.Increment("provider")

	// establish a global lock to ensure we're the only ones setting up a provider
	// also, we need to ensure the lock is released even if we error
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if err := Run("provider_setup", setup.control); err != nil {
		return err
	}

	return nil
}

// setupMounts will add the shares and mounts for this app
func (setup *devSetup) setupMounts() error {

	// mount the engine if it's a local directory
	if util.EngineDir() != "" {
		src := util.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), util.AppName())

		// first export the share on the workstation
		if err := provider.AddShare(src, dst); err != nil {
			return err
		}

		// mount the share on the provider
		if err := provider.AddMount(src, dst); err != nil {
			return err
		}
	}

	// mount the app src
	src := util.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), util.AppName())

	// first export the share on the workstation
	if err := provider.AddShare(src, dst); err != nil {
		return err
	}

	// then mount the share on the provider
	if err := provider.AddMount(src, dst); err != nil {
		return err
	}

	return nil
}

// setupApp sets up the app plaftorm and data services
func (setup *devSetup) setupApp() error {

	// let anyone else know we're using the app
	counter.Increment(util.AppName())

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	// setup the app
	if err := Run("app_setup", setup.control); err != nil {
		return err
	}

	// clean up after any possible failures in a previous deploy
	if err := Run("service_clean", setup.control); err != nil {
		return err
	}

	// setup the platform services
	return Run("platform_setup", setup.control)
}
