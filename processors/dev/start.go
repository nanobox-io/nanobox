package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
)

// Start initializes and starts the dev environment
func Start(envModel *models.Env, appModel *models.App) error {
	
	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}
	
	// setup the app if it's not already active
	switch appModel.State {
	case "active":
		// start the app
		if err := app.Start(appModel); err != nil {
			return fmt.Errorf("failed to start app: %s", err.Error())
		}		
	default:
		if err := app.Setup(envModel, appModel, "dev"); err != nil {
			return fmt.Errorf("failed to setup app: %s", err.Error())
		}	
	}

	return nil
}
