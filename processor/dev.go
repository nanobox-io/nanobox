package processor

import (
	"fmt"
	"io/ioutil"
	"os"
	"errors"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/counter"
)

type dev struct {
	config 				ProcessConfig
	oldBoxfile		models.Boxfile
	newBoxfile		models.Boxfile
	buildBoxfile 	models.Boxfile
}

func init() {
	Register("dev", devFunc)
}

func devFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["dev-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return dev{config: config}, nil
}

func (self dev) Results() ProcessConfig {
	return self.config
}

func (self dev) Process() error {

	if err := self.setup(); err != nil {
		// todo: how to display this?
		goto CLEANUP
	}

	if err := self.runDevSession(); err != nil {
		// todo: how to display this?
		goto CLEANUP
	}

CLEANUP:

	if err := self.teardown(); err != nil {
		// this is bad, really bad...
		// we should probably print a pretty message explaining that the app
		// was left in a bad state and needs to be reset
		os.Exit(1)
	}

	return nil
}

// setup sets up the provider, platform, and app
func (self *dev) setup() error {

	if err := self.setupProvider(); err != nil {
		return err
	}

	if err := self.setupApp(); err != nil {
		return err
	}

	return nil
}

// setupProvider sets up the provider
func (self *dev) setupProvider() error {

	// let anyone else know we're using the provider
	counter.Increment("provider")

	// establish a global lock to ensure we're the only ones setting up a provider
	// also, we need to ensure the lock is released even if we error
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if err := self.runProviderSetup(); err != nil {
		return err
	}

	return nil
}

// setupApp sets up the app plaftorm and data services
func (self *dev) setupApp() error {

	// let anyone else know we're using the app
	counter.Increment(util.AppName())

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	if err := self.runServiceClean(); err != nil {
		return err
	}

	if err := self.runPlatformSetup(); err != nil {
		return err
	}

	if err := self.fetchOldBoxfile(); err != nil {
		return err
	}

	if err := self.fetchNewBoxfile(); err != nil {
		return err
	}

	// todo: we need to consider a more stateful way of determining if a dev
	// was successful previously. Otherwise a failure on the first run won't
	// try a subsequent build
	if changed := self.hasBoxfileChanged(); changed == true {

		if err := self.persistNewBoxfile(); err != nil {
			return err
		}

		if err := self.runBuild(); err != nil {
			return err
		}

		if err := self.syncServices(); err != nil {
			return err
		}

	}

	if err := self.startDataServices(); err != nil {
		return err
	}

	return nil
}

// teardown tears down the app, platform, and provider
func (self *dev) teardown() error {

	if err := self.teardownApp(); err != nil {
		return err
	}

	if err := self.teardownProvider(); err != nil {
		return err
	}

	return nil
}

// teardownApp tears down the app when it's not being used
func (self *dev) teardownApp() error {

	counter.Decrement(util.AppName())

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	if unused := appIsUnused(); unused == true {

		if err := self.runPlatformStop(); err != nil {
			return err
		}

		if err := self.stopDataServices(); err != nil {
			return err
		}
	}

	return nil
}

// teardownProvider tears down the provider when it's not being used
func (self *dev) teardownProvider() error {

	counter.Decrement("provider")

	// establish a global lock to ensure we're the only ones bringing down
	// the provider. Also we need to ensure that we release the lock even
	// if we error out.
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if unused := providerIsUnused(); unused == true {
		if err := self.runProviderStop(); err != nil {
			return err
		}
	}

	return nil
}

// runProviderSetup sets up the docker environment
func (self *dev) runProviderSetup() error {
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		return err
	}

	return nil
}

// runProviderStop will stop the provider
func (self *dev) runProviderStop() error {
	err := Run("provider_stop", self.config)
	if err != nil {
		fmt.Println("provider_stop:", err)
		return err
	}

	return nil
}

// runServiceClean cleans up services that were left in a bad state
func (self *dev) runServiceClean() error {
	err := Run("service_clean", self.config)
	if err != nil {
		fmt.Println("service_clean:", err)
		return err
	}

	return nil
}

// runPlatformSetup sets up the platform services
func (self *dev) runPlatformSetup() error {
	err := Run("platform_setup", self.config)
	if err != nil {
		fmt.Println("platform_setup:", err)
		return err
	}

	return nil
}

// runPlatformStop stops the platform services
func (self *dev) runPlatformStop() error {
	err := Run("platform_stop", self.config)
	if err != nil {
		fmt.Println("platform_stop:", err)
		return err
	}

	return nil
}

// fetchOldBoxfile fetches the old boxfile from the db
func (self *dev) fetchOldBoxfile() error {
	// we don't care about the error here because it's very likely
	// that there won't be an old boxfile.
	data.Get(util.AppName()+"_meta", "boxfile", &self.oldBoxfile)

	return nil
}

// fetchNewBoxfile fetches the new boxfile
func (self *dev) fetchNewBoxfile() error {
	rawData, err := ioutil.ReadFile(util.BoxfileLocation())

	if err != nil {
		return errors.New("unable to load boxfile.yml")
	}

	self.newBoxfile.Data = rawData

	return nil
}

// persistNewBoxfile persists the new boxfile to the database
func (self *dev) persistNewBoxfile() error {

	key := util.AppName() + "_meta"
	if err := data.Put(key, "boxfile", self.newBoxfile); err != nil {
		return err
	}

	return nil
}

// hasBoxfileChanged returns true if the boxfile has changed
func (self *dev) hasBoxfileChanged() bool {

	if string(self.oldBoxfile.Data) != string(self.newBoxfile.Data) || len(self.oldBoxfile.Data) == 0 {
		return true
	}

	return false
}

// runBuild runs a build process
func (self *dev) runBuild() error {

	err := Run("code_build", self.config)
	if err != nil {
		fmt.Println("code_build:", err)
		return err
	}

	return nil
}

// syncServices syncs services to match the new boxfile
func (self *dev) syncServices() error {

	err := Run("service_sync", self.config)
	if err != nil {
		fmt.Println("service_sync:", err)
		return err
	}

	return nil
}

// startDataServices will start all data services
func (self *dev) startDataServices() error {
	err := Run("service_start_all", self.config)
	if err != nil {
		fmt.Println("service_start_all:", err)
	}

	return nil
}

// stopDataServices will stop all data services
func (self *dev) stopDataServices() error {
	err := Run("service_stop_all", self.config)
	if err != nil {
		fmt.Println("service_stop_all:", err)
	}

	return nil
}

// runDevSession starts a dev container and establishes a console session
func (self *dev) runDevSession() error {
	if err := Run("code_dev", self.config); err != nil {
		return err
	}

	return nil
}

// appIsUnused returns true if the app isn't being used by any other session
func appIsUnused() bool {
	count, err := counter.Get(util.AppName())

	if count == 0 && err == nil {
		return true
	}

	return false
}

// providerIsUnused returns true if the provider is currently not being used
func providerIsUnused() bool {
	count, err := counter.Get("provider")

	if count == 0 && err == nil {
		return true
	}

	return false
}
