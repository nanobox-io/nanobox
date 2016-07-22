package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
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

	// remove all dns entries for this app
	if err := processor.Run("env_dns_remove_all", teardown.control); err != nil {
		return err		
	}

	// teardown the mounts
	return teardown.teardownMounts()
}

// teardownApp tears down the app when it's not being used
func (teardown *processTeardown) teardownApp() error {

	// establish a local app lock to ensure we're the only ones bringing down the
	// app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// stop all data services
	if err := processor.Run("app_teardown", teardown.control); err != nil {
		return err
	}

	return nil
}

// teardownMounts removes the environments mounts from the provider
func (teardown *processTeardown) teardownMounts() error {
	// we will not be tearing down the mounts currently
	// this is because they are required for production builds
	return nil
	// return processor.Run("app_unmount", teardown.control)
}
