package service

import (
	"errors"
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processServiceStop ...
type processServiceStop struct {
	control processor.ProcessControl
	service models.Service
}

//
func init() {
	processor.Register("service_stop", serviceStopFunc)
}

//
func serviceStopFunc(control processor.ProcessControl) (processor.Processor, error) {
	if control.Meta["name"] == "" {
		return nil, errors.New("missing service name")
	}
	if control.Meta["label"] == "" {
		control.Meta["label"] = control.Meta["name"]
	}

	return processServiceStop{control: control}, nil
}

//
func (serviceStop processServiceStop) Results() processor.ProcessControl {
	return serviceStop.control
}

//
func (serviceStop processServiceStop) Process() error {
	if !serviceStop.isServiceRunning() {
		// short-circuit, this is already stopped
		return nil
	}

	if err := serviceStop.loadService(); err != nil {
		return err
	}

	if serviceStop.service.ID == "" {
		return errors.New("the service has not been created")
	}

	if err := serviceStop.stopContainer(); err != nil {
		return err
	}

	if err := serviceStop.detachNetwork(); err != nil {
		return err
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (serviceStop processServiceStop) isServiceRunning() bool {
	uid := serviceStop.control.Meta["name"]

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", config.AppName(), uid))

	return err == nil && container.State.Status == "running"
}

// loadService loads the service from the database
func (serviceStop *processServiceStop) loadService() error {
	// get the service from the database
	err := data.Get(config.AppName(), serviceStop.control.Meta["name"], &serviceStop.service)
	if err != nil {
		// cannot stop a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// stopContainer stops a docker container
func (serviceStop *processServiceStop) stopContainer() error {
	serviceStop.control.Display(stylish.Bullet("Stopping %s...", serviceStop.control.Meta["label"]))

	err := docker.ContainerStop(serviceStop.service.ID)
	if err != nil {
		return err
	}

	return nil
}

// detachNetwork detaches the container to the host network
func (serviceStop *processServiceStop) detachNetwork() error {

	if err := provider.RemoveNat(serviceStop.service.ExternalIP, serviceStop.service.InternalIP); err != nil {
		return err
	}

	if err := provider.RemoveIP(serviceStop.service.ExternalIP); err != nil {
		return err
	}

	return nil
}
