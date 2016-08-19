package platform

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
)

//
type Stop struct {
	App models.App
}

//
func (stop Stop) Run() error {

	// stop all the platform components weather they have been created
	// or not. (some may not be created in all instances)
	for _, component := range append(setupComponents, deployComponents...) {
		if err := stop.stopComponent(component); err != nil {
			return err
		}
	}

	return nil
}

// stopComponent will stop an individual component
func (stop *Stop) stopComponent(pComp Component) error {
	compModel, err := models.FindComponentBySlug(stop.App.ID, pComp.name)
	if err == nil {
		// if im able to retrieve the component from the db
		// stop it
		// probably put in some messaging here
		componentStop := component.Stop{compModel}
		return componentStop.Run()
	}

	return nil
}
