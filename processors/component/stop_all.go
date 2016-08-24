package component

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// StopAll stops all app components
func StopAll(a *models.App) error {
	// get all the components that belong to this app
	components, err := models.AllComponentsByApp(a.ID)
	if err != nil {
		return fmt.Errorf("unable to retrieve components: %s", err)
	}

	// stop each component
	for _, component := range components {
		if err := Stop(component); err != nil {
			return fmt.Errorf("unable to stop component(%s): %s", component.Name, err)
		}
	}

	return nil
}
