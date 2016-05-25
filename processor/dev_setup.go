package processor

import (
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/counter"
)

type devSetup struct {
	config 				ProcessConfig
}

func init() {
	Register("dev_setup", devSetupFunc)
}

func devSetupFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["devSetup-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return &devSetup{config: config}, nil
}

func (self devSetup) Results() ProcessConfig {
	return self.config
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

	return Run("provider_setup", self.config)
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
	if err := Run("service_clean", self.config); err != nil {
		return err
	}

	// setup the platform services
	return Run("platform_setup", self.config)
}