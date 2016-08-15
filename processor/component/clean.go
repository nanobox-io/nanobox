package component

import (
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
)

// Clean ...
type Clean struct {
	App models.App
}

//
func (clean Clean) Run() error {

	components, err := models.AllComponentsByApp(clean.App.ID)
	if err != nil {
		return err
	}

	for _, component := range components {
		if err := clean.cleanService(component); err != nil {
			return err
		}
	}

	return nil
}

// cleanService will clean a service if it was left in a bad state
func (clean Clean) cleanService(component models.Component) error {

	if clean.isComponentDirty(component) {
		return clean.removeService(component)
	}

	return nil
}

// removeService will remove a service from nanobox
func (clean Clean) removeService(component models.Component) error {

	componentRemove := Destroy{
		App:       clean.App,
		Component: component,
	}

	return componentRemove.Run()
}

// isComponentDirty will return true if the service is not active and available
func (clean Clean) isComponentDirty(component models.Component) bool {

	// short-circuit if this service never made it to active
	if component.State != ACTIVE {
		return true
	}

	return !clean.containerExists(component)
}

// containerExists will check to see if a docker container exists on the provider
func (clean Clean) containerExists(component models.Component) bool {
	_, err := docker.GetContainer(component.ID)
	return err == nil
}
