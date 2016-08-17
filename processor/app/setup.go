package app

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/processor/platform"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/locker"
)

//
type Setup struct {
	// required
	Env     models.Env
	AppName string

	// added
	App models.App
}

//
func (setup *Setup) Run() error {

	// fill in this apps name from the registry
	// this should allow the env.Setup to run
	// without ahving to knwo what apps will be setup
	if setup.AppName == "" {
		setup.AppName = registry.GetString("appname")
	}

	setup.loadApp()

	lumber.Debug("app load complete %+v", setup.App)

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is already active
	if setup.App.State == ACTIVE {
		return nil
	}

	if err := setup.reserveIPs(); err != nil {
		return err
	}

	if err := setup.generateEvars(); err != nil {
		return err
	}

	if err := setup.persistApp(); err != nil {
		return err
	}

	// clean up after any possible failures in a previous deploy
	componentClean := component.Clean{App: setup.App}
	if err := componentClean.Run(); err != nil {
		return err
	}

	// setup the platform services
	platformSetup := platform.Setup{App: setup.App}
	return platformSetup.Run()
}

// loadApp loads the app from the db
func (setup *Setup) loadApp() error {
	// the app might not exist yet, so let's not return the error if it fails
	setup.App, _ = models.FindAppBySlug(setup.Env.ID, setup.AppName)

	// set the default state
	if setup.App.State == "" {
		lumber.Debug("app not setup yet")
		setup.App.EnvID = setup.Env.ID
		setup.App.ID = fmt.Sprintf("%s_%s", setup.Env.ID, setup.AppName)
		setup.App.Name = setup.AppName
		setup.App.State = INITIALIZED
		setup.App.GlobalIPs = map[string]string{}
		setup.App.LocalIPs = map[string]string{}
		setup.App.Evars = map[string]string{}
	}

	lumber.Debug("app:Setup:loadApp:appModel:%+v", setup.App)

	return nil
}

// reserveIPs reserves necessary app global and local ip addresses
func (setup *Setup) reserveIPs() error {

	// reserve a dev ip
	envIP, err := dhcp.ReserveGlobal()
	if err != nil {
		lumber.Error("app:Setup:reserveIPs:dhcp.ReserveGlobal(): %s", err.Error())
		return err
	}

	// reserve a logvac ip
	logvacIP, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("app:Setup:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
		return err
	}

	// reserve a mist ip
	mistIP, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("app:Setup:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
		return err
	}

	// now let's assign them onto the app
	setup.App.GlobalIPs["env"] = envIP.String()

	setup.App.LocalIPs["logvac"] = logvacIP.String()
	setup.App.LocalIPs["mist"] = mistIP.String()

	return nil
}

// generateEvars generates the default app evars
func (setup *Setup) generateEvars() error {

	if setup.App.Evars["APP_NAME"] == "" {
		setup.App.Evars["APP_NAME"] = setup.AppName
	}

	return nil
}

// persistApp saves the app to the db
func (setup *Setup) persistApp() error {
	// set the app state to active so we don't setup again
	setup.App.State = ACTIVE
	
	if err := setup.App.Save(); err != nil {
		lumber.Error("app:Setup:persistApp:App.Save(): %s", err.Error())
		return err
	}
	return nil
}
