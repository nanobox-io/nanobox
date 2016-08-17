package sim

import (
	
	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/app"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"

)

// Start ...
type Start struct {
	// mandatory
	Env models.Env
	// created
	App models.App
}

//
func (start *Start) Run() error {
	display.OpenContext("starting sim")

	registry.Set("appname", "sim")

	if  start.App.ID == "" {
		start.App, _ = models.FindAppBySlug(config.EnvID(), "sim")
	}

	// if the app has not been setup
	// setup the app first
	if provider.Status() != "Running" {
		envSetup := env.Setup{}
		err := envSetup.Run()
		if err != nil {
			return err
		}
		start.Env = envSetup.Env
	} else {
		// only run the init if we dont do an env setup
		envInit := env.Init{}
		if err := envInit.Run(); err != nil {
			return err
		}
	}

	// retrieve the app
	a, _ := models.FindAppBySlug(start.Env.ID, "sim")

	// if the provider was running but the app wasnt setup
	if a.ID == "" {
		appSetup := app.Setup{Env: start.Env}
		if err := appSetup.Run(); err != nil {
			return err
		}
		a = appSetup.App
	}

	// run an app start
	appStart := &app.Start{App: a}
	if err := appStart.Run(); err != nil {
		return err
	}
	start.App = appStart.App

	display.CloseContext()
	return nil
}
