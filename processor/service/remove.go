package service

import (
	"errors"
	"net"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ipControl"
)

// processServiceRemove ...
type processServiceRemove struct {
	control processor.ProcessControl
	fail    bool
	service models.Service
}

//
func init() {
	processor.Register("service_remove", serviceRemoveFunc)
}

//
func serviceRemoveFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &processServiceRemove{control: control}, nil
}

//
func (serviceRemove processServiceRemove) Results() processor.ProcessControl {
	return serviceRemove.control
}

//
func (serviceRemove *processServiceRemove) Process() error {

	if err := serviceRemove.validateName(); err != nil {
		return err
	}

	if err := serviceRemove.loadService(); err != nil {
		// short-circuit if this service doesn't exist
		return nil
	}

	if err := serviceRemove.removeNat(); err != nil {
		return err
	}

	if err := serviceRemove.removeIP(); err != nil {
		return err
	}

	if err := serviceRemove.removeContainer(); err != nil {
		return err
	}

	if err := serviceRemove.releaseIPs(); err != nil {
		return err
	}

	if err := serviceRemove.deleteService(); err != nil {
		return err
	}

	return nil
}

// validateName validates a name was provided in the metadata
func (serviceRemove *processServiceRemove) validateName() error {
	if serviceRemove.control.Meta["name"] == "" {
		return errors.New("missing name")
	}

	return nil
}

// loadService loads the service from the database
func (serviceRemove *processServiceRemove) loadService() error {
	name := serviceRemove.control.Meta["name"]
	if err := data.Get(config.AppName(), name, &serviceRemove.service); err != nil {
		return err
	}

	return nil
}

// removeNat removes the nat from the container
func (serviceRemove *processServiceRemove) removeNat() error {
	extIP := serviceRemove.service.ExternalIP
	intIP := serviceRemove.service.InternalIP

	if err := provider.RemoveNat(extIP, intIP); err != nil {
		return err
	}

	return nil
}

// removeIP removes the IP from the host
func (serviceRemove *processServiceRemove) removeIP() error {
	if err := provider.RemoveIP(serviceRemove.service.ExternalIP); err != nil {
		return err
	}

	return nil
}

// releaseIPs releases the IPs back to the pool
func (serviceRemove *processServiceRemove) releaseIPs() error {
	extIP := net.ParseIP(serviceRemove.service.ExternalIP)
	intIP := net.ParseIP(serviceRemove.service.InternalIP)

	if err := ipControl.ReturnIP(intIP); err != nil {
		return err
	}

	if err := ipControl.ReturnIP(extIP); err != nil {
		return err
	}

	return nil
}

// removeContainer removes a container from the provider
func (serviceRemove *processServiceRemove) removeContainer() error {

	if exists := serviceRemove.containerExists(serviceRemove.service.ID); exists != true {
		return nil
	}

	if err := docker.ContainerRemove(serviceRemove.service.ID); err != nil {
		return err
	}

	return nil
}

// deleteServices removes the service entry from the database
func (serviceRemove *processServiceRemove) deleteService() error {

	name := serviceRemove.control.Meta["name"]
	if err := data.Delete(config.AppName(), name); err != nil {
		return err
	}

	return nil
}

// TODO: this should be a general helper in the docker library
// containerExists checks to see if a container exists
func (serviceRemove *processServiceRemove) containerExists(id string) bool {

	if _, err := docker.GetContainer(id); err == nil {
		return true
	}

	return false
}
