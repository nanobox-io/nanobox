package sim

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/app"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/util/config"
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

	start.App, _ = models.FindAppBySlug(config.EnvID(), "sim")

	// if the app has not been setup
	// setup the app first
	if start.App.ID == "" {
		envSetup := env.Setup{}
		err := envSetup.Run()
		if err != nil {
			return err
		}
	}

	// make sure we have an up to date app
	start.App, _ = models.FindAppBySlug(config.EnvID(), "sim")

	appSetup := app.Setup{App: start.App}
	if err := appSetup.Run(); err != nil {
		return err
	}

	// messaging about what happened and next steps

	return nil
}
