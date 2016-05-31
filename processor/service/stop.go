package service

import (
	"errors"
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStop struct {
	control  processor.ProcessControl
	service models.Service
}

func init() {
	processor.Register("service_stop", serviceStopFunc)
}

func serviceStopFunc(control processor.ProcessControl) (processor.Processor, error) {
	if control.Meta["name"] == "" {
		return nil, errors.New("missing service name")
	}	
	if control.Meta["label"] == "" {
		control.Meta["label"] = control.Meta["name"]
	}

	return serviceStop{control: control}, nil
}

func (self serviceStop) Results() processor.ProcessControl {
	return self.control
}

func (self serviceStop) Process() error {
	if !self.isServiceRunning() {
		// short-circuit, this is already stopped
		return nil
	}

	if err := self.loadService(); err != nil {
		return err
	}

	if self.service.ID == "" {
		return errors.New("the service has not been created")
	}

	if err := self.stopContainer(); err != nil {
		return err
	}

	if err := self.detachNetwork(); err != nil {
		return err
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (self serviceStop) isServiceRunning() bool {
	uid := self.control.Meta["name"]

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", util.AppName(), uid))

	return err == nil && container.State.Status == "running"
}

// loadService loads the service from the database
func (self *serviceStop) loadService() error {
	// get the service from the database
	err := data.Get(util.AppName(), self.control.Meta["name"], &self.service)
	if err != nil {
		// cannot stop a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// stopContainer stops a docker container
func (self *serviceStop) stopContainer() error {
	self.control.Display(stylish.Bullet("Stopping %s...", self.control.Meta["label"]))

	err := docker.ContainerStop(self.service.ID)
	if err != nil {
		return err
	}

	return nil
}

// detachNetwork detaches the container to the host network
func (self *serviceStop) detachNetwork() error {

	if err := provider.RemoveNat(self.service.ExternalIP, self.service.InternalIP); err != nil {
		return err
	}

	if err := provider.RemoveIP(self.service.ExternalIP); err != nil {
		return err
	}

	return nil
}
