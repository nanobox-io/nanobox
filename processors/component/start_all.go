package component

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

// StartAll starts all app components
func StartAll(a *models.App) error {
	// get all the components that belong to this app
	components, err := models.AllComponentsByApp(a.ID)
	if err != nil {
		lumber.Error("component:StartAll:models.AllComponentsByApp(%s): %s", a.ID, err.Error())
		return fmt.Errorf("unable to retrieve app components: %s", err)
	}

	// start each component
	for _, component := range components {
		if err := Start(component); err != nil {
			return fmt.Errorf("unable to start component(%s): %s", component.Name, err.Error())
		}
	}

	return nil
}
