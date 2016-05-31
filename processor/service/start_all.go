package service

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStartAll struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("service_start_all", serviceStartAllFunc)
}

func serviceStartAllFunc(control processor.ProcessControl) (processor.Processor, error) {
	// make sure i was given a name and image
	return serviceStartAll{control: control}, nil
}

func (self serviceStartAll) Results() processor.ProcessControl {
	return self.control
}

func (self serviceStartAll) Process() error {

	if err := self.startServices(); err != nil {
		return err
	}

	return nil
}

// startServices starts all of the services saved in the database
func (self serviceStartAll) startServices() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := self.startService(service); err != nil {
			return err
		}
	}

	return nil
}

// startService starts a service
func (self serviceStartAll) startService(uid string) error {

	config := processor.ProcessControl{
		DevMode: self.control.DevMode,
		Verbose: self.control.Verbose,
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
