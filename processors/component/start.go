package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Start starts the component services
func Start(componentModel *models.Component) error {

	// short-circuit if the container is already running
	if isComponentRunning(componentModel.ID) {
		return nil
	}

	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// make sure the component is active
	if componentModel.State != "active" {
		return fmt.Errorf("tried to start an inactive component")
	}

	// start the container
	if err := startContainer(componentModel.ID); err != nil {
		return err
	}

	return nil
}

// startContainer starts the container for this component
func startContainer(id string) error {
	display.StartTask("Start docker container")
	defer display.StopTask()

	if err := docker.ContainerStart(id); err != nil {
		lumber.Error("component:Start:docker.ContainerStart(%s): %s", id, err.Error())
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}

	return nil
}
