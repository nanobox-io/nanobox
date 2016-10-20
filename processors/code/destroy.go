package code

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

// Destroy destroys a code component from the app
func Destroy(componentModel *models.Component) error {
	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// remove the docker container
	if err := destroyContainer(componentModel.ID); err != nil {
		return err
	}

	// detach from the host network
	if err := detachNetwork(componentModel); err != nil {
		return fmt.Errorf("failed to detach container from the host network: %s", err.Error())
	}

	// remove the componentModel from the database
	if err := componentModel.Delete(); err != nil {
		lumber.Error("code:Destroy:Component.Delete(): %s", err.Error())
		display.ErrorTask()
		return fmt.Errorf("unable to delete database model: %s", err.Error())
	}

	return nil
}

// destroys a docker container associated with this app
func destroyContainer(id string) error {
	display.StartTask("Destroying docker container")
	defer display.StopTask()

	if id == "" {
		return nil
	}

	if err := docker.ContainerRemove(id); err != nil {
		lumber.Error("component:Destroy:docker.ContainerRemove(%s): %s", id, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to remove docker container: %s", err.Error())
	}

	return nil
}

// detachNetwork detaches the network from the host
func detachNetwork(componentModel *models.Component) error {
	display.StartTask("Releasing IPs")
	defer display.StopTask()

	//
	if err := provider.RemoveNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("code:Destroy:provider.RemoveNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("unable to remove network bridge: %s", err.Error())
	}

	//
	if err := provider.RemoveIP(componentModel.ExternalIP); err != nil {
		lumber.Error("code:Destroy:provider.RemoveIP(%s): %s", componentModel.ExternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("unable to release ip: %s", err.Error())
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(componentModel.ExternalIP)); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", componentModel.ExternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("unable to release ip: %s", err.Error())
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(componentModel.InternalIP)); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", componentModel.InternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf(": %s", err.Error())
	}

	return nil
}
