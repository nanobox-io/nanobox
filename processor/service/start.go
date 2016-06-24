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

// processServiceStart ...
type processServiceStart struct {
	control processor.ProcessControl
	service models.Service
}

//
func init() {
	processor.Register("service_start", serviceStartFn)
}

//
func serviceStartFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	// make sure i have a name to start
	if control.Meta["name"] == "" {
		return nil, errors.New("missing service name")
	}
	// set the label if it is missing
	if control.Meta["label"] == "" {
		control.Meta["label"] = control.Meta["name"]
	}

	return &processServiceStart{control: control}, nil
}

//
func (serviceStart processServiceStart) Results() processor.ProcessControl {
	return serviceStart.control
}

//
func (serviceStart *processServiceStart) Process() error {

	if running := serviceStart.isServiceRunning(); running == true {
		// short-circuit, this is already running
		return nil
	}

	if err := serviceStart.loadService(); err != nil {
		return err
	}

	if serviceStart.service.State != ACTIVE {
		return errors.New("the service has not been created")
	}

	if err := serviceStart.startContainer(); err != nil {
		return err
	}

	if err := serviceStart.attachNetwork(); err != nil {
		return err
	}

	return nil
}

// loadService loads the service from the database
func (serviceStart *processServiceStart) loadService() error {
	bucket := fmt.Sprintf("%s_%s", config.AppName(), serviceStart.control.Env)

	// get the service from the database
	err := data.Get(bucket, serviceStart.control.Meta["name"], &serviceStart.service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// startContainer starts a docker container
func (serviceStart *processServiceStart) startContainer() error {
	header := fmt.Sprintf("Starting %s...", serviceStart.control.Meta["label"])
	serviceStart.control.Info(stylish.Bullet(header))

	err := docker.ContainerStart(serviceStart.service.ID)
	if err != nil {
		return err
	}

	return nil
}

// attachNetwork attaches the container to the host network
func (serviceStart *processServiceStart) attachNetwork() error {
	fmt.Println("service start setting up network")
	// todo: add these to a cleanup process in case of failure

	err := provider.AddIP(serviceStart.service.ExternalIP)
	if err != nil {
		return err
	}

	err = provider.AddNat(serviceStart.service.ExternalIP, serviceStart.service.InternalIP)
	if err != nil {
		return err
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (serviceStart processServiceStart) isServiceRunning() bool {
	uid := serviceStart.control.Meta["name"]

	container, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s-%s", config.AppName(), serviceStart.control.Env, uid))

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
