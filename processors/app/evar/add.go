package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/display"
)

func Add(envModel *models.Env, appModel *models.App, evars map[string]string) error {

	if err := app.Setup(envModel, appModel, appModel.Name); err != nil {
		return fmt.Errorf("failed to setup app: %s", err)
	}

	// iterate through the evars and add them to the app
	for key, val := range evars {
		appModel.Evars[key] = val
	}

	// save the app
	if err := appModel.Save(); err != nil {
		return fmt.Errorf("failed to persist evars: %s", err.Error())
	}

	// iterate one more time for display
	fmt.Println()
	for key := range evars {
		fmt.Printf("%s %s added\n", display.TaskComplete, key)
	}
	fmt.Println()

	return nil
}
