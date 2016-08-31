package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
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

	if err := startNetwork(componentModel); err != nil {
		return err
	}

	// todo: set status

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

// startNetwork attaches the container to the virtual network
func startNetwork(componentModel *models.Component) error {
	display.StartTask("Attaching network")
	defer display.StopTask()
	
	// add the IP to the provider
	if err := provider.AddIP(componentModel.ExternalIP); err != nil {
		display.ErrorTask()
		lumber.Error("component:Start:provider.AddIP(%s): %s", componentModel.ExternalIP, err.Error())
		return fmt.Errorf("failed to add IP to the provider: %s", err.Error())
	}

	// nat traffic to the container
	if err := provider.AddNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		display.ErrorTask()
		lumber.Error("component:Start:attachNetwork:provider.AddNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		return fmt.Errorf("failed to setup NAT on the provider: %s", err.Error())
	}
	
	return nil
}
