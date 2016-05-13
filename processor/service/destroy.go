package service

import (
	"net"
	"errors"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type serviceDestroy struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("service_destroy", serviceDestroyFunc)
}

func serviceDestroyFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	if config.Meta["name"] == "" {
		return nil, errors.New("missing image or name")
	}
	return &serviceDestroy{config: config}, nil
}

func (self serviceDestroy) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceDestroy) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &service)
	if err != nil {
		// cant find service
		return err
	}

	err = docker.ContainerRemove(service.ID)
	if err != nil {
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

	err = ip_control.ReturnIP(net.ParseIP(service.ExternalIP))
	if err != nil {
		return err
	}

	err = ip_control.ReturnIP(net.ParseIP(service.InternalIP))
	if err != nil {
		return err
	}

	// save the service
	return data.Delete(util.AppName(), self.config.Meta["name"])
}
