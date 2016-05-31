package service

import (
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStopAll struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("service_stop_all", serviceStopAllFunc)
}

func serviceStopAllFunc(control processor.ProcessControl) (processor.Processor, error) {
	return serviceStopAll{control: control}, nil
}

func (self serviceStopAll) Results() processor.ProcessControl {
	return self.control
}

func (self serviceStopAll) Process() error {

	self.control.Display(stylish.Bullet("Stopping All Services"))

	if err := self.stopServices(); err != nil {
		return err
	}

	return nil
}

// stopServices stops all of the services saved in the database
func (self serviceStopAll) stopServices() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := self.stopService(service); err != nil {
			return err
		}
	}

	return nil
}

// stopService stops a service
func (self serviceStopAll) stopService(uid string) error {

	config := processor.ProcessControl{
		DevMode: self.control.DevMode,
		Verbose: self.control.Verbose,
		DisplayLevel: self.control.DisplayLevel+1,
		Meta: map[string]string{
			"name":  uid,
		},
	}

	// provision
	if err := processor.Run("service_stop", config); err != nil {
		return err
	}

	return nil
}
