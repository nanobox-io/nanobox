package service

import (
	"errors"
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStart struct {
	control  processor.ProcessControl
	service models.Service
}

func init() {
	processor.Register("service_start", serviceStartFunc)
}

func serviceStartFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	// make sure i have a name to start
	if control.Meta["name"] == "" {
		return nil, errors.New("missing service name")
	}
	// set the label if it is missing
	if control.Meta["label"] == "" {
		control.Meta["label"] = control.Meta["name"]
	}

	return &serviceStart{control: control}, nil
}

func (self serviceStart) Results() processor.ProcessControl {
	return self.control
}

func (self *serviceStart) Process() error {

	if running := self.isServiceRunning(); running == true {
		// short-circuit, this is already running
		return nil
	}

	if err := self.loadService(); err != nil {
		return err
	}

	if self.service.State != "active" {
		return errors.New("the service has not been created")
	}

	if err := self.startContainer(); err != nil {
		return err
	}

	if err := self.attachNetwork(); err != nil {
		return err
	}

	return nil
}

// loadService loads the service from the database
func (self *serviceStart) loadService() error {
	// get the service from the database
	err := data.Get(util.AppName(), self.control.Meta["name"], &self.service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// startContainer starts a docker container
func (self *serviceStart) startContainer() error {
	header := fmt.Sprintf("Starting %s...", self.control.Meta["label"])
	self.control.Info(stylish.Bullet(header))

	err := docker.ContainerStart(self.service.ID)
	if err != nil {
		return err
	}

	return nil
}

// attachNetwork attaches the container to the host network
func (self *serviceStart) attachNetwork() error {

	// todo: add these to a cleanup process in case of failure

	err := provider.AddIP(self.service.ExternalIP)
	if err != nil {
		return err
	}

	err = provider.AddNat(self.service.ExternalIP, self.service.InternalIP)
	if err != nil {
		return err
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (self serviceStart) isServiceRunning() bool {
	uid := self.control.Meta["name"]

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", util.AppName(), uid))

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
