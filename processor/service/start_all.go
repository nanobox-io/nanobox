package service

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processServiceStartAll ...
type processServiceStartAll struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("service_start_all", serviceStartAllFn)
}

//
func serviceStartAllFn(control processor.ProcessControl) (processor.Processor, error) {
	// make sure i was given a name and image
	return processServiceStartAll{control: control}, nil
}

//
func (serviceStartAll processServiceStartAll) Results() processor.ProcessControl {
	return serviceStartAll.control
}

//
func (serviceStartAll processServiceStartAll) Process() error {

	if err := serviceStartAll.startServices(); err != nil {
		return err
	}

	return nil
}

// startServices starts all of the services saved in the database
func (serviceStartAll processServiceStartAll) startServices() error {
	bucket := fmt.Sprintf("%s_%s", config.AppName(), serviceStartAll.control.Env)
	
	services, err := data.Keys(bucket)
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := serviceStartAll.startService(service); err != nil {
			return err
		}
	}

	return nil
}

// startService starts a service
func (serviceStartAll processServiceStartAll) startService(uid string) error {

	config := processor.ProcessControl{
		Env: serviceStartAll.control.Env,
		Verbose: serviceStartAll.control.Verbose,
		Meta: map[string]string{
			"label": uid,
			"name":  uid,
		},
	}

	// provision
	if err := processor.Run("service_start", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_start:", uid), err)
		return err
	}

	return nil
}
