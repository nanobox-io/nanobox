package component

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
)

// Start ...
type Start struct {
	Component models.Component
}

//
func (start *Start) Run() error {

	// short-circuit if the service is running
	if start.isServiceRunning() {
		return nil
	}

	if start.Component.State != ACTIVE {
		return fmt.Errorf("the service has not been created")
	}

	if err := start.startContainer(); err != nil {
		return err
	}

	if err := start.attachNetwork(); err != nil {
		return err
	}

	return nil
}

// startContainer starts a docker container
func (start *Start) startContainer() error {

	err := docker.ContainerStart(start.Component.ID)
	if err != nil {
		return err
	}

	return nil
}

// attachNetwork attaches the container to the host network
func (start *Start) attachNetwork() error {
	err := provider.AddIP(start.Component.ExternalIP)
	if err != nil {
		return err
	}

	err = provider.AddNat(start.Component.ExternalIP, start.Component.InternalIP)
	if err != nil {
		return err
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (start Start) isServiceRunning() bool {
	container, err := docker.GetContainer(start.Component.ID)

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
