package code

import (
	"net"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

// Destroy ...
type Destroy struct {
	Component models.Component
}

//
func (destroy *Destroy) Run() error {

	//
	if err := docker.ContainerRemove(destroy.Component.ID); err != nil {
		lumber.Error("code:Destroy:docker.ContainerRemove(%s): %s", destroy.Component.ID, err.Error())
		return err
	}

	//
	if err := provider.RemoveNat(destroy.Component.ExternalIP, destroy.Component.InternalIP); err != nil {
		lumber.Error("code:Destroy:provider.RemoveNat(%s, %s): %s", destroy.Component.ExternalIP, destroy.Component.InternalIP, err.Error())
		return err
	}

	//
	if err := provider.RemoveIP(destroy.Component.ExternalIP); err != nil {
		lumber.Error("code:Destroy:provider.RemoveIP(%s): %s", destroy.Component.ExternalIP, err.Error())
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(destroy.Component.ExternalIP)); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", destroy.Component.ExternalIP, err.Error())
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(destroy.Component.InternalIP)); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", destroy.Component.InternalIP, err.Error())
		return err
	}

	// remove the destroy.Component from the database
	if err := destroy.Component.Delete(); err != nil {
		lumber.Error("code:Destroy:Component.Delete(): %s", err.Error())
		return err
	}

	return nil
}
