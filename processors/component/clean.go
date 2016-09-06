package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Clean purges any components in a dirty or incomplete state
func Clean(appModel *models.App) error {
	// fetch all of the app components
	components, err := appModel.Components()
	if err != nil {
		lumber.Error("component:Clean:models.App{ID:%s}.Components(): %s", appModel.ID, err.Error())
		return fmt.Errorf("failed to fetch app component collection: %s", err.Error())
	}

	if !areComponentsDirty(components) {
		return nil
	}

	display.OpenContext("Cleaning dirty components")
	defer display.CloseContext()

	// iterate through the components and clean them
	for _, componentModel := range components {
		if err := cleanComponent(appModel, componentModel); err != nil {
			return fmt.Errorf("failed to clean component: %s", err.Error())
		}
	}

	return nil
}

// cleanComponent will clean a component if it was left in a bad state
func cleanComponent(appModel *models.App, componentModel *models.Component) error {

	// short-circuit if the component is not dirty
	if !isComponentDirty(componentModel) {
		return nil
	}

	if err := Destroy(appModel, componentModel); err != nil {
		return fmt.Errorf("failed to remove component: %s", err.Error())
	}

	return nil
}

// areComponentsDirty checks to see if any of the components are dirty
func areComponentsDirty(componentModels []*models.Component) bool {
	for _, componentModel := range componentModels {
		if isComponentDirty(componentModel) {
			return true
		}
	}

	return false
}

// isComponentDirty returns true if the container is removed or in a bad state
func isComponentDirty(componentModel *models.Component) bool {
	// short-circuit if this service never made it to active
	if componentModel.State != "active" {
		return true
	}

	// check to see if the container exists
	_, err := docker.GetContainer(componentModel.ID)
	return err != nil
}
