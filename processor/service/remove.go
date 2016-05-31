package service

import (
	"errors"
	"net"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type serviceRemove struct {
	control  processor.ProcessControl
	fail    bool
	service models.Service
}

func init() {
	processor.Register("service_remove", serviceRemoveFunc)
}

func serviceRemoveFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &serviceRemove{control: control}, nil
}

func (self serviceRemove) Results() processor.ProcessControl {
	return self.control
}

func (self *serviceRemove) Process() error {

	if err := self.validateName(); err != nil {
		return err
	}

	if err := self.loadService(); err != nil {
		// short-circuit if this service doesn't exist
		return nil
	}

	if err := self.removeNat(); err != nil {
		return err
	}

	if err := self.removeIP(); err != nil {
		return err
	}

	if err := self.removeContainer(); err != nil {
		return err
	}

	if err := self.releaseIPs(); err != nil {
		return err
	}

	if err := self.deleteService(); err != nil {
		return err
	}

	return nil
}

// validateName validates a name was provided in the metadata
func (self *serviceRemove) validateName() error {
	if self.control.Meta["name"] == "" {
		return errors.New("missing name")
	}

	return nil
}

// loadService loads the service from the database
func (self *serviceRemove) loadService() error {
	name := self.control.Meta["name"]
	if err := data.Get(util.AppName(), name, &self.service); err != nil {
		return err
	}

	return nil
}

// removeNat removes the nat from the container
func (self *serviceRemove) removeNat() error {
	extIP := self.service.ExternalIP
	intIP := self.service.InternalIP

	if err := provider.RemoveNat(extIP, intIP); err != nil {
		return err
	}

	return nil
}

// removeIP removes the IP from the host
func (self *serviceRemove) removeIP() error {
	if err := provider.RemoveIP(self.service.ExternalIP); err != nil {
		return err
	}

	return nil
}

// releaseIPs releases the IPs back to the pool
func (self *serviceRemove) releaseIPs() error {
	extIP := net.ParseIP(self.service.ExternalIP)
	intIP := net.ParseIP(self.service.InternalIP)

	if err := ip_control.ReturnIP(intIP); err != nil {
		return err
	}

	if err := ip_control.ReturnIP(extIP); err != nil {
		return err
	}

	return nil
}

// removeContainer removes a container from the provider
func (self *serviceRemove) removeContainer() error {

	if exists := self.containerExists(self.service.ID); exists != true {
		return nil
	}

	if err := docker.ContainerRemove(self.service.ID); err != nil {
		return err
	}

	return nil
}

// deleteServices removes the service entry from the database
func (self *serviceRemove) deleteService() error {

	name := self.control.Meta["name"]
	if err := data.Delete(util.AppName(), name); err != nil {
		return err
	}

	return nil
}

// todo: this should be a general helper in the docker library
// containerExists checks to see if a container exists
func (self *serviceRemove) containerExists(id string) bool {

	if _, err := docker.GetContainer(id); err == nil {
		return true
	}

	return false
}
