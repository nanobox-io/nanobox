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
	serviceStartAll := &processServiceStartAll{control: control}
	return serviceStartAll, serviceStartAll.validateMeta()
}

//
func (serviceStartAll processServiceStartAll) Results() processor.ProcessControl {
	return serviceStartAll.control
}

//
func (serviceStartAll *processServiceStartAll) Process() error {

	// get the service keys
	services, err := data.Keys(serviceStartAll.control.Meta["name"])
	if err != nil {
		return err
	}

	// start each service
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
			"app_name": serviceStartAll.control.Meta["app_name"],
			"label": uid,
			"name":  uid,
		},
	}

	// provision
	if err := processor.Run("service_start", config); err != nil {
		// serviceStartAll.control.Display(fmt.Sprintf("%s_start: %+v", uid, err))
		return err
	}

	return nil
}

// validateMeta validates the meta data
// it also sets a default for the name of the app
func (serviceStartAll *processServiceStartAll) validateMeta() error {

	// set the name of the app if we are not given one
	if serviceStartAll.control.Meta["app_name"] == "" {
		serviceStartAll.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppName(), serviceStartAll.control.Env)
	}

	return nil
}
