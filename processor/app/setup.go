package app

import (
	"net"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
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
	data.Get("apps", config.AppName(), &appSetup.app)

	// set the default state
	if appSetup.app.State == "" {
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
	var devIP net.IP
	devIP, err = dhcp.ReserveGlobal()
	if err != nil {
		return err
	}

	// reserve a preview ip
	var previewIP net.IP
	previewIP, err = dhcp.ReserveGlobal()
	if err != nil {
		return err
	}

	// reserve a logvac ip
	var logvacIP net.IP
	logvacIP, err = dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	// reserve a mist ip
	var mistIP net.IP
	mistIP, err = dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	// now let's assign them onto the app
	appSetup.app.GlobalIPs["dev"] 		= devIP.String()
	appSetup.app.GlobalIPs["preview"] = previewIP.String()

	appSetup.app.LocalIPs["logvac"] = logvacIP.String()
	appSetup.app.LocalIPs["mist"]   = mistIP.String()

	return nil
}

// generateEvars generates the default app evars
func (appSetup *processAppSetup) generateEvars() error {
	// fetch the app evars model if it exists
	evars := models.EnvVars{}

	// ignore the error because it's likely to not exist
	data.Get(config.AppName()+"_meta", "env", &evars)

	if evars["APP_NAME"] == "" {
		evars["APP_NAME"] = config.AppName()
	}

	if err := data.Put(config.AppName()+"_meta", "env", evars); err != nil {
		return err
	}

	return nil
}

// persistApp saves the app to the db
func (appSetup *processAppSetup) persistApp() error {

	// set the app state to active so we don't appSetup again
	appSetup.app.State = ACTIVE

	// save the app
	if err := data.Put("apps", config.AppName(), &appSetup.app); err != nil {
		return err
	}

	return nil
}
