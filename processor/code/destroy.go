package code

import (
	"net"

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
		return err
	}

	//
	if err := provider.RemoveNat(destroy.Component.ExternalIP, destroy.Component.InternalIP); err != nil {
		return err
	}

	//
	if err := provider.RemoveIP(destroy.Component.ExternalIP); err != nil {
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(destroy.Component.ExternalIP)); err != nil {
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(destroy.Component.InternalIP)); err != nil {
		return err
	}

	// remove the destroy.Component from the database
	return destroy.Component.Delete()
}
