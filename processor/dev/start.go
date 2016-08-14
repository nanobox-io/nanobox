package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/processor/app"
	"github.com/nanobox-io/nanobox/processor/env"
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
	registry.Set("appname", "dev")

	envInit := env.Init{}
	if err := envInit.Run(); err != nil {
		return err
	}

	// if the app has not been setup
	// setup the app first
	if !start.appExists(start.Env.ID) {
		envSetup := env.Setup{}
		err := envSetup.Run()
		if err != nil {
			return err
		}
	}

	start.App, _ = models.FindAppBySlug(start.Env.ID, "dev")
	appStart := app.Start{App: start.App}
	if err := appStart.Run(); err != nil {
		return err
	}

	// messaging about what happened and next steps

	return nil
}

func (start *Start) appExists(envID string) bool {
	appModel, _ := models.FindAppBySlug(start.Env.ID, "dev")
	return appModel.ID != ""
}

