package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Start starts the component services
func Start(c *models.Component) error {

	// short-circuit if the container is already running
	if isComponentRunning(c.ID) {
		return nil
	}

	display.OpenContext(c.Label)
	defer display.CloseContext()

	// make sure the component is active
	if c.State != "active" {
		return fmt.Errorf("tried to start an inactive component")
	}

	// start the container
	if err := docker.ContainerStart(c.ID); err != nil {
		lumber.Error("component:Start:docker.ContainerStart(%s): %s", c.ID, err.Error())
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}

	// add the IP to the provider
	if err := provider.AddIP(c.ExternalIP); err != nil {
		lumber.Error("component:Start:provider.AddIP(%s): %s", c.ExternalIP, err.Error())
		return fmt.Errorf("failed to add IP to the provider: %s", err.Error())
	}

	// nat traffic to the container
	if err := provider.AddNat(c.ExternalIP, c.InternalIP); err != nil {
		lumber.Error("component:Start:attachNetwork:provider.AddNat(%s, %s): %s", c.ExternalIP, c.InternalIP, err.Error())
		return fmt.Errorf("failed to setup NAT on the provider: %s", err.Error())
	}

	// todo: set status

	return nil
}
