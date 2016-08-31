package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
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
	
	// stop from the virtual network
	if err := stopNetwork(componentModel); err != nil {
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

// stopNetwork stops the network on this component
func stopNetwork(componentModel *models.Component) error {
	display.StartTask("Detaching network")
	defer display.StopTask()
	
	// remove NAT
	if err := provider.RemoveNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		display.ErrorTask()
		lumber.Error("component:stopNetwork:provider.RemoveNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		return fmt.Errorf("failed to remove NAT on the provider: %s", err.Error())
	}

	// remove the IP from the provider
	if err := provider.RemoveIP(componentModel.ExternalIP); err != nil {
		display.ErrorTask()
		lumber.Error("component:stopNetwork:provider.RemoveIP(%s): %s", componentModel.ExternalIP, err.Error())
		return fmt.Errorf("failed to remove IP from the provider: %s", err.Error())
	}
	
	return nil
}
