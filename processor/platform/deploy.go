package platform

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
)

// Deploy ...
type Deploy struct {
	App models.App
}

//
func (deploy Deploy) Run() error {

	for _, component := range deployComponents {
		if err := deploy.provisionComponent(component); err != nil {
			return err
		}
	}

	return nil
}

// provisionComponent will provision an individual component
func (deploy Deploy) provisionComponent(pComp Component) error {

	// if the component exists and is active just start it and stop here
	if deploy.isComponentActive(pComp.name) {
		compModel, _ := models.FindComponentBySlug(deploy.App.ID, pComp.name)
		componentStart := component.Start{Component: compModel}
		return componentStart.Run()
	}

	// otherwise
	// deploy the component
	componentSetup := component.Setup{
		App:   deploy.App,
		Name:  pComp.name,
		Image: pComp.image,
	}
	if err := componentSetup.Run(); err != nil {
		return err
	}

	// and configure it
	componentConfigure := component.Configure{
		App:       deploy.App,
		Component: componentSetup.Component,
	}
	return componentConfigure.Run()
}

// isComponentActive returns true if a component is already active
func (deploy Deploy) isComponentActive(name string) bool {

	// component db entry
	component, _ := models.FindComponentBySlug(deploy.App.ID, name)

	return component.State == "active"
}
