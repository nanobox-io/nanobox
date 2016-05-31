package processor

import (
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

func (self devSetup) Results() ProcessControl {
	return self.control
}

func (self *devSetup) Process() error {
	if err := self.setupProvider(); err != nil {
		return err
	}

	return self.setupApp()
}

// setupProvider sets up the provider
func (self *devSetup) setupProvider() error {

	// let anyone else know we're using the provider
	counter.Increment("provider")

	// establish a global lock to ensure we're the only ones setting up a provider
	// also, we need to ensure the lock is released even if we error
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	return Run("provider_setup", self.control)
}

// setupApp sets up the app plaftorm and data services
func (self *devSetup) setupApp() error {

	// let anyone else know we're using the app
	counter.Increment(util.AppName())

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	// clean up after any possible failures in a previous deploy
	if err := Run("service_clean", self.control); err != nil {
		return err
	}

	// setup the platform services
	return Run("platform_setup", self.control)
}
