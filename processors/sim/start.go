package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
)

// Start initializes and starts the sim environment
func Start(envModel *models.Env, appModel *models.App) error {

	// initialize the environemnt
	if err := env.Setup(envModel); err != nil {
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
		if err := app.Setup(envModel, appModel, "sim"); err != nil {
			return fmt.Errorf("failed to setup app: %s", err.Error())
		}
	}

	return nil
}
