package component

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

// StopAll stops all app components
func StopAll(appModel *models.App) error {

	// get all the components that belong to this app
	componentModels, err := appModel.Components()
	if err != nil {
		lumber.Error("component:StopAll:models.App{ID:%s}.Components() %s", appModel.ID, err.Error())
		return fmt.Errorf("unable to retrieve components: %s", err.Error())
	}

	// stop each component
	for _, componentModel := range componentModels {
		if err := Stop(componentModel); err != nil {
			return fmt.Errorf("unable to stop component(%s): %s", componentModel.Name, err.Error())
		}
	}

	return nil
}
