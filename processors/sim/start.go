package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Start initializes and starts the dev environment
func Start(envModel *models.Env, appModel *models.App) error {

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
		if err := app.Setup(envModel, appModel); err != nil {
			return fmt.Errorf("failed to setup app: %s", err.Error())
		}
	}

	// start the app
	if err := app.Start(envModel, appModel); err != nil {
		return fmt.Errorf("failed to start app: %s", err.Error())
	}

	return nil
}
