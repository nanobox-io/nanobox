package dev

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Start initializes and starts the dev environment
func Start(e *models.Env, a *models.App) error {

	// ensure the provider is setup
	if err := provider.Setup(); err != nil {
		return fmt.Errorf("failed to setup the provider: %s", err.Error())
	}
	
	// ensure the environment is setup
	if err := env.Setup(e); err != nil {
		return fmt.Errorf("failed to setup the env: %s", err.Error())
	}

	// setup the app if it's not already active
	if a.State != "active" {
		if err := app.Setup(e, a); err != nil {
			return fmt.Errorf("failed to setup app: %s", err.Error())
		}
	}

	// start the app
	if err := app.Start(e, a); err != nil {
		return fmt.Errorf("failed to start app: %s", err.Error())
	}

	return nil
}
