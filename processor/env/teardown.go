package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processTeardown ...
type processTeardown struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("env_teardown", teardownFn)
}

//
func teardownFn(control processor.ProcessControl) (processor.Processor, error) {
	return &processTeardown{control}, nil
}

//
func (teardown processTeardown) Results() processor.ProcessControl {
	return teardown.control
}

//
func (teardown *processTeardown) Process() error {

	// dont shut anything down if we are supposed to background
	if processor.DefaultControl.Debug {
		fmt.Println("leaving running because your running in debug mode")
		return nil
	}

	// if im given a environment to teardown i will tear the app down 
	// in the env. If not (in the case of a build) no app is teardownable.
	if teardown.control.Env != "" {
		if err := teardown.teardownApp(); err != nil {
			return err
		}
	}

	if err := teardown.teardownMounts(); err != nil {
		return err
	}

	if err := teardown.teardownProvider(); err != nil {
		return err
	}

	return nil
}

// teardownApp tears down the app when it's not being used
func (teardown *processTeardown) teardownApp() error {

	counter.Decrement(config.AppName())

	// establish a local app lock to ensure we're the only ones bringing down the
	// app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is still in use
	if appInUse() {
		return nil
	}

	// Stop the platform services
	if err := processor.Run("platform_stop", teardown.control); err != nil {
		return err
	}

	// stop all data services
	if err := processor.Run("service_stop_all", teardown.control); err != nil {
		return err
	}

	return nil
}

// teardownMounts removes the unused mounts from the provider
func (teardown *processTeardown) teardownMounts() error {

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is still in use
	if appInUse() {
		return nil
	}

	// unmount the engine if it's a local directory
	if config.EngineDir() != "" {

		//
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

		// unmount the env on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			return err
		}
	}

	// unmount the app src
	src := config.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.AppName())

	// unmount the env on the provider
	if err := provider.RemoveMount(src, dst); err != nil {
		return err
	}

	return nil
}

// teardownProvider tears down the provider when it's not being used
func (teardown *processTeardown) teardownProvider() error {

	//
	count, err := counter.Decrement("provider")
	if err != nil {
		return err
	}

	// establish a global lock to ensure we're the only ones bringing down
	// the provider. Also we need to ensure that we release the lock even
	// if we error out.
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	// stop the provider
	if providerIsUnused() {
		return processor.Run("provider_stop", teardown.control)
	}

	//
	teardown.control.Display(stylish.Bullet("the provider is needed for %d more thing(s)", count))

	return nil
}

// appInUse returns true if the app is being used by any other session
func appInUse() bool {
	count, err := counter.Get(config.AppName())
	return err != nil || count != 0
}

// providerIsUnused returns true if the provider is currently not being used
func providerIsUnused() bool {
	count, err := counter.Get("provider")
	return err == nil && count == 0
}
