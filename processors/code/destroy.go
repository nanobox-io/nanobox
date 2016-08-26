package code

import (
	"net"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Destroy ...
//
func Destroy(componentModel *models.Component) error {

	//
	if err := docker.ContainerRemove(componentModel.ID); err != nil {
		lumber.Error("code:Destroy:docker.ContainerRemove(%s): %s", componentModel.ID, err.Error())
		return err
	}

	//
	if err := provider.RemoveNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("code:Destroy:provider.RemoveNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		return err
	}

	//
	if err := provider.RemoveIP(componentModel.ExternalIP); err != nil {
		lumber.Error("code:Destroy:provider.RemoveIP(%s): %s", componentModel.ExternalIP, err.Error())
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(componentModel.ExternalIP)); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", componentModel.ExternalIP, err.Error())
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(componentModel.InternalIP)); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", componentModel.InternalIP, err.Error())
		return err
	}

	// remove the componentModel from the database
	if err := componentModel.Delete(); err != nil {
		lumber.Error("code:Destroy:Component.Delete(): %s", err.Error())
		return err
	}

	return nil
}
