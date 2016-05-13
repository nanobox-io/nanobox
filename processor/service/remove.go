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

type serviceRemove struct {
	config processor.ProcessConfig
	fail   bool
}

func init() {
	processor.Register("service_remove", serviceRemoveFunc)
}

func serviceRemoveFunc(config processor.ProcessConfig) (processor.Processor, error) {
	return &serviceRemove{config: config}, nil
}

func (self serviceRemove) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceRemove) Process() error {
	// make sure i was given a name and image
	if self.config.Meta["name"] == "" {
		return errors.New("missing image or name")
	}

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &service)
	if err != nil {
		// quit early if the service was found to be created already
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
	err = docker.ContainerRemove(service.ID)
	if err != nil {
		return err
	}
	err = ip_control.ReturnIP(net.ParseIP(service.InternalIP))
	if err != nil {
		return err
	}
	err = ip_control.ReturnIP(net.ParseIP(service.ExternalIP))
	if err != nil {
		return err
	}
	return data.Delete(util.AppName(), self.config.Meta["name"])
}
