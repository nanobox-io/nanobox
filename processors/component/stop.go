package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Stop stops the component's docker container
func Stop(componentModel *models.Component) error {
	// short-circuit if the process is already stopped
	if !isComponentRunning(componentModel.ID) {
		return nil
	}

	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// stop the docker container
	if err := stopContainer(componentModel.ID); err != nil {
		return err
	}

	return nil
}

// stopContainer stops the docker container for this component
func stopContainer(id string) error {
	display.StartTask("Stopping docker container")
	defer display.StopTask()

	if err := docker.ContainerStop(id); err != nil {
		display.ErrorTask()
		lumber.Error("component:Stop:docker.ContainerStop(%s): %s", id, err.Error())
		return fmt.Errorf("failed to stop docker container: %s", err.Error())
	}

	return nil
}
