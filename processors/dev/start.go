package dev

import (
	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

type Start struct {
	// mandatory
	Env models.Env
	// created
	App models.App
}

// Run initializes and starts the dev environment
func (start *Start) Run() error {
	display.OpenContext("starting dev")

	// before we can setup the app, the env has to be setup
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
	a, _ := models.FindAppBySlug(start.Env.ID, "dev")

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
