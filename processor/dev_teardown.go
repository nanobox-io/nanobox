package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processDevTeardown ...
type processDevTeardown struct {
	control ProcessControl
}

//
func init() {
	Register("dev_teardown", devTeardownFunc)
}

//
func devTeardownFunc(control ProcessControl) (Processor, error) {
	return &processDevTeardown{control}, nil
}

//
func (devTeardown processDevTeardown) Results() ProcessControl {
	return devTeardown.control
}

//
func (devTeardown *processDevTeardown) Process() error {
	// dont shut anything down if we are supposed to background
	if DefaultConfig.Background {
		return nil
	}

	if err := devTeardown.teardownApp(); err != nil {
		return err
	}

	if err := devTeardown.teardownMounts(); err != nil {
		return err
	}

	if err := devTeardown.teardownProvider(); err != nil {
		return err
	}

	return nil
}

// teardownApp tears down the app when it's not being used
func (devTeardown *processDevTeardown) teardownApp() error {

	counter.Decrement(config.AppName())

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	if appIsUnused() {

		// Stop the platform services
		if err := Run("platform_stop", devTeardown.control); err != nil {
			return err
		}

		// stop all data services
		if err := Run("service_stop_all", devTeardown.control); err != nil {
			return err
		}
	}

	return nil
}

// teardownMounts removes the unused mounts from the provider
func (devTeardown *processDevTeardown) teardownMounts() error {

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is still in use
	if !appIsUnused() {
		return nil
	}

	// unmount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

		// unmount the share on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
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

	return nil
}

// teardownProvider tears down the provider when it's not being used
func (devTeardown *processDevTeardown) teardownProvider() error {

	count, err := counter.Decrement("provider")
	if err != nil {
		return err
	}

	// establish a global lock to ensure we're the only ones bringing down
	// the provider. Also we need to ensure that we release the lock even
	// if we error out.
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if providerIsUnused() {
		// stop the provider
		return Run("provider_stop", devTeardown.control)
	}
	devTeardown.control.Display(stylish.Bullet("%d dev's still running so leaving the provider up", count))
	return nil
}

// appIsUnused returns true if the app isn't being used by any other session
func appIsUnused() bool {
	count, err := counter.Get(config.AppName())
	return err == nil && count == 0
}

// providerIsUnused returns true if the provider is currently not being used
func providerIsUnused() bool {
	count, err := counter.Get("provider")
	return err == nil && count == 0
}
