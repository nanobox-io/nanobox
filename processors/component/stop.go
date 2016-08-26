package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Stop stops the component's docker container
func Stop(componentModel *models.Component) error {
	// short-circuit if the process is already stopped
	if !isComponentRunning(componentModel.ID) {
		return nil
	}

	// stop the docker container
	if err := docker.ContainerStop(componentModel.ID); err != nil {
		lumber.Error("component:Stop:docker.ContainerStop(%s): %s", componentModel.ID, err.Error())
		return fmt.Errorf("failed to stop docker container: %s", err.Error())
	}

	// remove NAT
	if err := provider.RemoveNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("component:Stop:provider.RemoveNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		return fmt.Errorf("failed to remove NAT on the provider: %s", err.Error())
	}

	// remove the IP from the provider
	if err := provider.RemoveIP(componentModel.ExternalIP); err != nil {
		lumber.Error("component:Stop:provider.RemoveIP(%s): %s", componentModel.ExternalIP, err.Error())
		return fmt.Errorf("failed to remove IP from the provider: %s", err.Error())
	}

	return nil
}
