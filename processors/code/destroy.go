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

// Destroy ...
//
func Destroy(componentModel *models.Component) error {
	display.OpenContext("removing %s(%s) component", componentModel.Label, componentModel.Name)
	defer display.CloseContext()

	//
	display.StartTask("removing container")
	if err := docker.ContainerRemove(componentModel.ID); err != nil {
		lumber.Error("code:Destroy:docker.ContainerRemove(%s): %s", componentModel.ID, err.Error())
		display.ErrorTask()
		return err
	}
	display.StopTask()

	display.StartTask("cleaning networking")
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

	// remove the componentModel from the database
	if err := componentModel.Delete(); err != nil {
		lumber.Error("code:Destroy:Component.Delete(): %s", err.Error())
		display.ErrorTask()
		return fmt.Errorf("unable to delete database model: %s", err.Error())
	}

	display.StopTask()

	return nil
}
