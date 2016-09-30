package component

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Destroy destroys a component from the provider and database
func Destroy(appModel *models.App, componentModel *models.Component) error {
	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// remove the docker container
	if err := destroyContainer(componentModel.ID); err != nil {
		return err
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

	// remove NAT
	if err := provider.RemoveNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("component:detachNetwork:provider.RemoveNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to remove NAT from provider: %s", err.Error())
	}

	// remove IP
	if err := provider.RemoveIP(componentModel.ExternalIP); err != nil {
		lumber.Error("component:detachNetwork:provider.RemoveIP(%s): %s", componentModel.ExternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to remove IP from provider: %s", err.Error())
	}

	// return the external IP
	// don't return the external IP if this is portal
	if componentModel.Name != "portal" && appModel.GlobalIPs[componentModel.Name] == "" {
		ip := net.ParseIP(componentModel.ExternalIP)
		if err := dhcp.ReturnIP(ip); err != nil {
			lumber.Error("component:detachNetwork:dhcp.ReturnIP(%s): %s", ip, err.Error())
			display.ErrorTask()
			return fmt.Errorf("failed to release IP back to pool: %s", err.Error())
		}
	}

	// return the internal IP
	// don't return the internal IP if it's an app-level cache
	if appModel.LocalIPs[componentModel.Name] == "" {
		ip := net.ParseIP(componentModel.InternalIP)
		if err := dhcp.ReturnIP(ip); err != nil {
			lumber.Error("component:detachNetwork:dhcp.ReturnIP(%s): %s", ip, err.Error())
			display.ErrorTask()
			return fmt.Errorf("failed to release IP back to pool: %s", err.Error())
		}
	}

	return nil
}
