package component

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// StartAll ...
type StartAll struct {
	App models.App
}

//
func (startAll *StartAll) Run() error {
	// get all the components that belong to this app
	components, err := models.AllComponentsByApp(startAll.App.ID)
	if err != nil {
		return fmt.Errorf("unable to retrieve components: %s", err)
	}

	// start each component
	for _, component := range components {
		if err := startAll.startComponent(component); err != nil {
			return fmt.Errorf("unable to start component(%s): %s", component.Name, err)
		}
	}

	return nil
}

// startComponent starts a component
func (startAll StartAll) startComponent(component models.Component) error {
	start := Start{component}
	return start.Run()
}
