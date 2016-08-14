package platform

import (

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
)

// Setup ...
type Setup struct {
	App models.App
}

//
func (setup Setup) Run() error {

	for _, component := range setupComponents {
		if err := setup.provisionComponent(component); err != nil {
			return err
		}
	}

	return nil
}


// provisionComponent will provision an individual component
func (setup Setup) provisionComponent(pComp Component) error {

	if setup.isComponentActive(pComp.name) {
		// start the component if the component is already active
		comp, _  := models.FindComponentBySlug(setup.App.ID, pComp.name)
		componentStart := component.Start{Component: comp}
		return componentStart.Run()
	}

	// otherwise
	// setup the component
	componentSetup := component.Setup{
		App: setup.App,
		Name: pComp.name,
		Image: pComp.image,
	}
	if err := componentSetup.Run(); err != nil {
		return err
	}

	// and configure it
	componentConfigure := component.Configure{
		App: setup.App,
		Component: componentSetup.Component,
	}
	return componentConfigure.Run()
}

// isComponentActive returns true if a component is already active
func (setup Setup) isComponentActive(name string) bool {

	// component db entry
	component, _  := models.FindComponentBySlug(setup.App.ID, name)

	return component.State == "active"
}
