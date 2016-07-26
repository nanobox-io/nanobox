package service

import (
	"fmt"

	"github.com/jcelliott/lumber"
	
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// processServiceStop ...
type processServiceStop struct {
	control processor.ProcessControl
	service models.Service
	label   string
	name    string
}

//
func init() {
	processor.Register("service_stop", serviceStopFn)
}

//
func serviceStopFn(control processor.ProcessControl) (processor.Processor, error) {
	serviceStop := &processServiceStop{control: control}
	return serviceStop, serviceStop.validateMeta()
}

//
func (serviceStop *processServiceStop) Results() processor.ProcessControl {
	return serviceStop.control
}

//
func (serviceStop *processServiceStop) Process() error {

	// short-circuit if the process is already stopped
	if !serviceStop.isServiceRunning() {
		return nil
	}

	// attempt to load the service
	if err := serviceStop.loadService(); err != nil {
		return err
	}

	// attempt to stop the container
	if err := serviceStop.stopContainer(); err != nil {
		return err
	}

	// attempt to detach the network
	if err := serviceStop.detachNetwork(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (serviceStop *processServiceStop) validateMeta() error {

	// set meta values
	serviceStop.name = serviceStop.control.Meta["name"]
	serviceStop.label = serviceStop.control.Meta["label"]

	// ensure name is provided
	if serviceStop.name == "" {
		return fmt.Errorf("Missing meta data 'name'")
	}

	// if no label is provided, just use the name
	if serviceStop.label == "" {
		serviceStop.label = serviceStop.name
	}

	// set the name of the app if we are not given one
	if serviceStop.control.Meta["app_name"] == "" {
		serviceStop.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppID(), serviceStop.control.Env)
	}

	return nil
}

// isServiceRunning returns true if a service is already running
func (serviceStop *processServiceStop) isServiceRunning() bool {

	// get the container
	container, err := docker.GetContainer(fmt.Sprintf("nanobox_%s_%s", serviceStop.control.Meta["app_name"], serviceStop.name))

	if err != nil {
		// we cant return an error but we can definatly log what happened
		lumber.Error("Service Stop I failed to retrieve nanobox_%s_%s_%s\n%s", config.AppID(), serviceStop.control.Env, serviceStop.name, err.Error())
		return false
	}

	return container.State.Status == "running"
}

// loadService loads the service from the database; an error here means we cannot
// find a service in the database
func (serviceStop *processServiceStop) loadService() error {

	//
	if err := data.Get(serviceStop.control.Meta["app_name"], serviceStop.name, &serviceStop.service); err != nil {
		print.OutputProcessorErr("service not found", fmt.Sprintf(`
Nanobox was unable to find the service '%s' in the database, try shutting
down again.
		`, serviceStop.name))

		return err
	}

	//
	if serviceStop.service.ID == "" {
		print.OutputProcessorErr("service not created", fmt.Sprintf(`
Nanobox was unable to stop '%s' because it has not yet been created
		`, serviceStop.name))

		return fmt.Errorf("Service not created")
	}

	return nil
}

// stopContainer stops a docker container
func (serviceStop *processServiceStop) stopContainer() error {

	serviceStop.control.Display(stylish.Bullet("Stopping %s...", serviceStop.label))

	//
	if err := docker.ContainerStop(serviceStop.service.ID); err != nil {
		print.OutputProcessorErr("failed to stop container", fmt.Sprintf(`
Nanobox failed to stop the container '%s', try shutting down again.
		`, serviceStop.service.ID))
	}

	return nil
}

// detachNetwork detaches the container from the host network
func (serviceStop *processServiceStop) detachNetwork() error {

	//
	if err := provider.RemoveNat(serviceStop.service.ExternalIP, serviceStop.service.InternalIP); err != nil {
		print.OutputProcessorErr("failed to remove nat", fmt.Sprintf(`
Nanobox failed to remove the NAT, try shutting down again.
		`))

		return err
	}

	//
	if err := provider.RemoveIP(serviceStop.service.ExternalIP); err != nil {
		print.OutputProcessorErr("failed to remove ip", fmt.Sprintf(`
Nanobox failed to remove the external IP, try shutting down again.
		`))

		return err
	}

	return nil
}
