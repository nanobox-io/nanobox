package service

import (
	"errors"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStop struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("service_stop", serviceStopFunc)
}

func serviceStopFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// make sure i was given a name and image
	if config.Meta["name"] == "" {
		return nil, errors.New("missing image or name")
	}

	return serviceStop{config: config}, nil
}

func (self serviceStop) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceStop) Process() error {
	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	err = provider.RemoveNat(service.ExternalIP, service.InternalIP)
	if err != nil {
		return err
	}

	err = provider.RemoveIP(service.ExternalIP)
	if err != nil {
		return err
	}

	return docker.ContainerStop(service.ID)
}
