package component

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
)

// Destroy destroys a component from the provider and database
func Destroy(appModel *models.App, componentModel *models.Component) error {
	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// remove the docker container
	if err := destroyContainer(componentModel.ID); err != nil {
		// report the error but continue on
		lumber.Error("component:Destroy:destroyContainer(%s): %s", componentModel.ID, err)
		// return err
	}

	// detach from the host network
	if err := detachNetwork(appModel, componentModel); err != nil {
		return fmt.Errorf("failed to detach container from the host network: %s", err.Error())
	}

	// purge evars
	if err := componentModel.PurgeEvars(appModel); err != nil {
		lumber.Error("component:Destroy:models.Component.PurgeEvars(%+v): %s", appModel, err.Error())
		return fmt.Errorf("failed to purge component evars from app: %s", err.Error())
	}

	// destroy the data model
	if err := componentModel.Delete(); err != nil {
		lumber.Error("component:Destroy:models.Component.Delete(): %s", err.Error())
		return fmt.Errorf("failed to destroy component model: %s", err.Error())
	}

	return nil
}

// destroyContainer destroys a docker container associated with this component
func destroyContainer(id string) error {
	display.StartTask("Destroying docker container")
	defer display.StopTask()

	// if i dont know the id then i cant remove it
	if id == "" {
		return nil
	}

	if err := docker.ContainerRemove(id); err != nil {
		lumber.Error("component:Destroy:docker.ContainerRemove(%s): %s", id, err.Error())
		display.ErrorTask()
		// return fmt.Errorf("failed to remove docker container: %s", err.Error())
	}

	return nil
}

// detachNetwork detaches the network from the host
func detachNetwork(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Releasing IPs")
	defer display.StopTask()

	if componentModel.IPAddr() == "" {
		return nil
	}

	// return the external IP
	// don't return the external IP if this is portal
	if componentModel.Name != "portal" && appModel.LocalIPs[componentModel.Name] == "" {
		ip := net.ParseIP(componentModel.IPAddr())
		if err := dhcp.ReturnIP(ip); err != nil {
			lumber.Error("component:detachNetwork:dhcp.ReturnIP(%s): %s", ip, err.Error())
			display.ErrorTask()
			return fmt.Errorf("failed to release IP back to pool: %s", err.Error())
		}
	}

	return nil
}
