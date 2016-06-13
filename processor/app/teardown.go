package app

import (
	"net"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ipControl"
)

// processAppTeardown ...
type processAppTeardown struct {
	control processor.ProcessControl
	app     models.App
}

//
func init() {
	processor.Register("app_appTeardown", appTeardownFunc)
}

//
func appTeardownFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &processAppTeardown{control: control}, nil
}

//
func (appTeardown *processAppTeardown) Results() processor.ProcessControl {
	return appTeardown.control
}

//
func (appTeardown *processAppTeardown) Process() error {

	if err := appTeardown.loadApp(); err != nil {
		return err
	}

	// short-circuit if the app isn't active
	if appTeardown.app.State == "initialized" {
		return nil
	}

	if err := appTeardown.releaseIPs(); err != nil {
		return err
	}

	if err := appTeardown.deleteEvars(); err != nil {
		return err
	}

	if err := appTeardown.deleteApp(); err != nil {
		return err
	}

	return nil
}

// loadApp loads the app from the db
func (appTeardown *processAppTeardown) loadApp() error {
	// the app might not exist yet, so let's not return the error if it fails
	data.Get("apps", util.AppName(), &appTeardown.app)

	// set the default state
	if appTeardown.app.State == "" {
		appTeardown.app.State = "initialized"
	}

	return nil
}

// releaseIPs releases necessary app-global ip addresses
func (appTeardown *processAppTeardown) releaseIPs() error {

	if err := ipControl.ReturnIP(net.ParseIP(appTeardown.app.DevIP)); err != nil {
		return err
	}

	return nil
}

// deleteEvars deletes the evars from the db
func (appTeardown *processAppTeardown) deleteEvars() error {

	// delete the evars model
	if err := data.Delete(util.AppName()+"_meta", "env"); err != nil {
		return err
	}

	return nil
}

// deleteApp deletes the app to the db
func (appTeardown *processAppTeardown) deleteApp() error {

	// delete the app model
	if err := data.Delete("apps", util.AppName()); err != nil {
		return err
	}

	return nil
}
