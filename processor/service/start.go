package service

import (
	"fmt"
	"errors"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStart struct {
	config processor.ProcessConfig
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
	// make sure i was given a name and image
	if self.config.Meta["name"] == "" {
		return missingImageOrName
	}

	fmt.Println("-> starting", self.config.Meta["name"])
	// get the service from the database
	service := models.Service{}
	data.Get(util.AppName(), self.config.Meta["name"], &service)

	// create docker container
	if service.ID == "" {
		// quit early if the service was found to be created already
		return errors.New("the service has not been created")
	}

	err := docker.ContainerStart(service.ID)
	if err != nil {
		return err
	}

	err = provider.AddIP(service.ExternalIP)
	if err != nil {
		return err
	}

	err = provider.AddNat(service.ExternalIP, service.InternalIP)
	if err != nil {
		return err
	}
	return nil
}
