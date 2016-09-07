package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

func Add(appModel *models.App, evars map[string]string) error {

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
	for key, _ := range evars {
		fmt.Printf("%s %s added\n", display.TaskComplete, key)
	}
	fmt.Println()

	return nil
}
