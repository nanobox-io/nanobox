package processor

import (
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/counter"
)

type devTeardown struct {
	config ProcessConfig
}

func init() {
	Register("dev_teardown", devTeardownFunc)
}

func devTeardownFunc(config ProcessConfig) (Processor, error) {
	return &devTeardown{config}, nil
}

func (self devTeardown) Results() ProcessConfig {
	return self.config
}

func (self *devTeardown) Process() error {
	// dont shut anything down if we are supposed to background
	if DefaultConfig.Background {
		return nil
	}

	if err := self.teardownApp(); err != nil {
		return err
	}

	return self.teardownProvider()
}


// teardownApp tears down the app when it's not being used
func (self *devTeardown) teardownApp() error {

	counter.Decrement(util.AppName())

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	if unused := appIsUnused(); unused == true {

		// Stop the platform services
		if err := Run("platform_stop", self.config); err != nil {
			return err
		}

		// stop all data services
		if err := Run("service_stop_all", self.config); err != nil {
			return err
		}
	}

	return nil
}

// teardownProvider tears down the provider when it's not being used
func (self *devTeardown) teardownProvider() error {

	counter.Decrement("provider")

	// establish a global lock to ensure we're the only ones bringing down
	// the provider. Also we need to ensure that we release the lock even
	// if we error out.
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if unused := providerIsUnused(); unused == true {
		// stop the provider
		if err := Run("provider_stop", self.config); err != nil {
			return err
		}
	}
	return nil
}

// appIsUnused returns true if the app isn't being used by any other session
func appIsUnused() bool {
	count, err := counter.Get(util.AppName())
	return count == 0 && err == nil
}

// providerIsUnused returns true if the provider is currently not being used
func providerIsUnused() bool {
	count, err := counter.Get("provider")
	return count == 0 && err == nil
}
