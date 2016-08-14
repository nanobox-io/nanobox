package component

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// StopAll ...
type StopAll struct {
	App models.App
}

//
func (stopAll *StopAll) Run() error {
	// get all the components that belong to this app
	components, err := models.AllComponentsByApp(stopAll.App.ID)
	if err != nil {
		return fmt.Errorf("unable to retrieve components: %s", err)
	}

	// stop each component
	for _, component := range components {
		if err := stopAll.stopComponent(component); err != nil {
			return fmt.Errorf("unable to stop component(%s): %s", component.Name, err)
		}
	}

	return nil
}

// stopComponent stops a service
func (stopAll *StopAll) stopComponent(component models.Component) error {
	stop := Stop{component}
	//
	return stop.Run()
}