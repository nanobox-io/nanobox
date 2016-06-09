package app

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type appSetup struct {
	control processor.ProcessControl
	app     models.App
}

func init() {
	processor.Register("app_setup", appSetupFunc)
}

func appSetupFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &appSetup{control: control}, nil
}

func (setup *appSetup) Results() processor.ProcessControl {
	return setup.control
}

func (setup *appSetup) Process() error {

	if err := setup.loadApp(); err != nil {
		return err
	}

	// short-circuit if the app is already active
	if setup.app.State == "active" {
		return nil
	}

	if err := setup.reserveIPs(); err != nil {
		return err
	}

	if err := setup.persistApp(); err != nil {
		return err
	}

	return nil
}

// loadApp loads the app from the db
func (setup *appSetup) loadApp() error {
	// the app might not exist yet, so let's not return the error if it fails
	data.Get("apps", util.AppName(), &setup.app)

	// set the default state
	if setup.app.State == "" {
		setup.app.State = "initialized"
	}

	return nil
}

// reserveIPs reserves necessary app-global ip addresses
func (setup *appSetup) reserveIPs() error {

	//
	devIP, err := ip_control.ReserveGlobal()
	if err != nil {
		return err
	}

	//
	setup.app.DevIP = devIP.String()

	return nil
}

// persistApp saves the app to the db
func (setup *appSetup) persistApp() error {

	// set the app state to active so we don't setup again
	setup.app.State = "active"

	// save the app
	if err := data.Put("apps", util.AppName(), &setup.app); err != nil {
		return err
	}

	return nil
}
