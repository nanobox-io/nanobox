package app

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

//
type processAppSetup struct {
	control processor.ProcessControl
	app     models.App
}

//
func init() {
	processor.Register("app_setup", appSetupFn)
}

//
func appSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	return &processAppSetup{control: control}, nil
}

//
func (appSetup *processAppSetup) Results() processor.ProcessControl {
	return appSetup.control
}

//
func (appSetup *processAppSetup) Process() error {

  // establish an app-level lock to ensure we're the only ones setting up an app
  // also, we need to ensure that the lock is released even if we error out.
  locker.LocalLock()
  defer locker.LocalUnlock()

	if err := appSetup.loadApp(); err != nil {
		return err
	}

	// short-circuit if the app is already active
	if appSetup.app.State == ACTIVE {
		return nil
	}

	if err := appSetup.reserveIPs(); err != nil {
		return err
	}

	if err := appSetup.generateEvars(); err != nil {
		return err
	}

	if err := appSetup.persistApp(); err != nil {
		return err
	}

	return nil
}

// loadApp loads the app from the db
func (appSetup *processAppSetup) loadApp() error {
	// the app might not exist yet, so let's not return the error if it fails
	key := fmt.Sprintf("%s_%s", config.AppName(), appSetup.control.Env)
	data.Get("apps", key, &appSetup.app)

	// set the default state
	if appSetup.app.State == "" {
		appSetup.app.Name = key
		appSetup.app.Directory = config.LocalDir()
		appSetup.app.State = INITIALIZED
		appSetup.app.GlobalIPs = map[string]string{}
		appSetup.app.LocalIPs = map[string]string{}
	}

	return nil
}

// reserveIPs reserves necessary app global and local ip addresses
func (appSetup *processAppSetup) reserveIPs() error {

	var err error

	// reserve a dev ip
	envIP, err := dhcp.ReserveGlobal()
	if err != nil {
		return err
	}

	// reserve a logvac ip
	logvacIP, err := dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	// reserve a mist ip
	mistIP, err := dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	// now let's assign them onto the app
	appSetup.app.GlobalIPs["env"] = envIP.String()

	appSetup.app.LocalIPs["logvac"] = logvacIP.String()
	appSetup.app.LocalIPs["mist"]   = mistIP.String()

	return nil
}

// generateEvars generates the default app evars
func (appSetup *processAppSetup) generateEvars() error {
	// create the bucket name one time
	bucket := fmt.Sprintf("%s_meta", config.AppName())

	// fetch the app evars model if it exists
	evars := models.Evars{}

	// ignore the error because it's likely to not exist
	data.Get(bucket, appSetup.control.Env+"_env", &evars)

	if evars["APP_NAME"] == "" {
		evars["APP_NAME"] = config.AppName()
	}

	if err := data.Put(bucket, appSetup.control.Env+"_env", evars); err != nil {
		return err
	}

	return nil
}

// persistApp saves the app to the db
func (appSetup *processAppSetup) persistApp() error {

	// set the app state to active so we don't appSetup again
	appSetup.app.State = ACTIVE

	// save the app
	key := fmt.Sprintf("%s_%s", config.AppName(), appSetup.control.Env)
	if err := data.Put("apps", key, &appSetup.app); err != nil {
		return err
	}

	return nil
}
