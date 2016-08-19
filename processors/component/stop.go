package component

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Stop ...
type Stop struct {
	Component models.Component
}

//
func (stop *Stop) Run() error {

	// short-circuit if the process is already stopped
	if !stop.isServiceRunning() {
		return nil
	}

	// attempt to stop the container
	if err := stop.stopContainer(); err != nil {
		return err
	}

	// attempt to detach the network
	if err := stop.detachNetwork(); err != nil {
		return err
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (stop *Stop) isServiceRunning() bool {

	// get the container
	container, err := docker.GetContainer(stop.Component.ID)

	if err != nil {
		// we cant return an error but we can definatly log what happened
		lumber.Error("Service Stop I failed to retrieve nanobox_%s_%s\n%s", stop.Component.AppID, stop.Component.Name, err.Error())
		return false
	}

	return container.State.Status == "running"
}

// stopContainer stops a docker container
func (stop *Stop) stopContainer() error {

	if err := docker.ContainerStop(stop.Component.ID); err != nil {
		lumber.Error("component:Stop:stopContainer:docker.ContainerStop(%s): %s", stop.Component.ID, err.Error())
		// TODO: display some error message but do not quit here
	}

	return nil
}

// detachNetwork detaches the container from the host network
func (stop *Stop) detachNetwork() error {

	//
	if err := provider.RemoveNat(stop.Component.ExternalIP, stop.Component.InternalIP); err != nil {
		lumber.Error("component:Stop:attachNetwork:provider.AddNat(%s, %s): %s", stop.Component.ExternalIP, stop.Component.InternalIP, err.Error())
		return err
	}

	//
	if err := provider.RemoveIP(stop.Component.ExternalIP); err != nil {
		lumber.Error("component:Stop:attachNetwork:provider.AddIP(%s): %s", stop.Component.ExternalIP, err.Error())
		return err
	}

	return nil
}
