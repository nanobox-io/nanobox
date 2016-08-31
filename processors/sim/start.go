package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
)

// Start initializes and starts the sim environment
func Start(envModel *models.Env, appModel *models.App) error {
	display.OpenContext("starting sim environemnt")
	defer display.CloseContext()

	// ensure the provider is setup
	if err := provider.Setup(); err != nil {
		return fmt.Errorf("failed to setup the provider: %s", err.Error())
	}

	// ensure the environment is setup
	if err := env.Setup(envModel); err != nil {
		return fmt.Errorf("failed to setup the env: %s", err.Error())
	}

	// setup the app if it's not already active
	if appModel.State != "active" {
		if err := app.Setup(envModel, appModel, "sim"); err != nil {
			return fmt.Errorf("failed to setup app: %s", err.Error())
		}
	}

	// start the app
	if err := app.Start(appModel); err != nil {
		return fmt.Errorf("failed to start app: %s", err.Error())
	}

	return nil
}
