package service

import (
	"fmt"
	"errors"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStop struct {
	config processor.ProcessConfig
	service 	models.Service
}

func init() {
	processor.Register("service_stop", serviceStopFunc)
}

func serviceStopFunc(config processor.ProcessConfig) (processor.Processor, error) {
	return serviceStop{config: config}, nil
}

func (self serviceStop) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceStop) Process() error {

	if err := self.validateMeta(); err != nil {
		return err
	}

	if running := self.isServiceRunning(); running == false {
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

// validateMeta validates that the provided metadata is supplied
func (self serviceStop) validateMeta() error {

  if self.config.Meta["label"] == "" {
    return errors.New("missing service label")
  }

  if self.config.Meta["name"] == "" {
    return errors.New("missing service name")
  }

  return nil
}

// isServiceRunning returns true if a service is already running
func (self serviceStop) isServiceRunning() bool {
	uid := self.config.Meta["name"]
	name := fmt.Sprintf("%s-%s", util.AppName(), uid)

	container, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	if err != nil {
		return false
	}

	// return true if the container is running
	if container.State.Status == "running" {
		return true
	}

	return false
}

// loadService loads the service from the database
func (self *serviceStop) loadService() error {
	// get the service from the database
	err := data.Get(util.AppName(), self.config.Meta["name"], &self.service)
	if err != nil {
		// cannot stop a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// stopContainer stops a docker container
func (self *serviceStop) stopContainer() error {
	header := fmt.Sprintf("Stopping %s...", self.config.Meta["Label"])
	fmt.Print(stylish.NestedBullet(header, self.config.DisplayLevel))

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
