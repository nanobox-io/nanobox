package component

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// StartAll starts all app components
func StartAll(a *models.App) error {

	// get all the components that belong to this app
	components, err := models.AllComponentsByApp(a.ID)
	if err != nil {
		lumber.Error("component:StartAll:models.AllComponentsByApp(%s): %s", a.ID, err.Error())
		return util.ErrorAppend(err, "unable to retrieve app components")
	}

	if len(components) == 0 {
		return nil
	}

	display.OpenContext("Starting components")
	defer display.CloseContext()

	// start each component
	for _, component := range components {
		if err := Start(component); err != nil {
			return util.ErrorAppend(err, "unable to start component(%s)", component.Name)
		}
	}

	return nil
}
