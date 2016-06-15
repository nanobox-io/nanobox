package service

import (
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processServiceStopAll ...
type processServiceStopAll struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("service_stop_all", serviceStopAllFn)
}

//
func serviceStopAllFn(control processor.ProcessControl) (processor.Processor, error) {
	return processServiceStopAll{control: control}, nil
}

//
func (serviceStopAll processServiceStopAll) Results() processor.ProcessControl {
	return serviceStopAll.control
}

//
func (serviceStopAll processServiceStopAll) Process() error {

	serviceStopAll.control.Display(stylish.Bullet("Stopping All Services"))

	if err := serviceStopAll.stopServices(); err != nil {
		return err
	}

	return nil
}

// stopServices stops all of the services saved in the database
func (serviceStopAll processServiceStopAll) stopServices() error {
	services, err := data.Keys(config.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := serviceStopAll.stopService(service); err != nil {
			return err
		}
	}

	return nil
}

// stopService stops a service
func (serviceStopAll processServiceStopAll) stopService(uid string) error {

	config := processor.ProcessControl{
		DevMode:      serviceStopAll.control.DevMode,
		Verbose:      serviceStopAll.control.Verbose,
		DisplayLevel: serviceStopAll.control.DisplayLevel + 1,
		Meta: map[string]string{
			"name": uid,
		},
	}

	// provision
	if err := processor.Run("service_stop", config); err != nil {
		return err
	}

	return nil
}
