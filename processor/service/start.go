package service

import (
	"fmt"
	"errors"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStart struct {
	config 		processor.ProcessConfig
	service 	models.Service
}

func init() {
	processor.Register("service_start", serviceStartFunc)
}

func serviceStartFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &serviceStart{config: config}, nil
}

func (self serviceStart) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceStart) Process() error {

	if err := self.validateImage(); err != nil {
		return err
	}

	if err := self.loadService(); err != nil {
		return err
	}

	if self.service.ID == "" {
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

// validateImage ensures an image or name was provided
func (self *serviceStart) validateImage() error {
	// make sure i was given a name and image
	if self.config.Meta["name"] == "" {
		return errors.New("missing image or name")
	}

	return nil
}

// loadService loads the service from the database
func (self *serviceStart) loadService() error {
	// get the service from the database
	err := data.Get(util.AppName(), self.config.Meta["name"], &self.service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// startContainer starts a docker container
func (self *serviceStart) startContainer() error {
	label := "Starting container..."
	fmt.Print(stylish.NestedProcessStart(label, self.config.DisplayLevel))

	err := docker.ContainerStart(self.service.ID)
	if err != nil {
		return err
	}

	return nil
}

// attachNetwork attaches the container to the host network
func (self *serviceStart) attachNetwork() error {
	label := "Add container to host network..."
	fmt.Print(stylish.NestedProcessStart(label, self.config.DisplayLevel))

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
