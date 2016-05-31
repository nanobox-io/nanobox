package platform

import (
	"github.com/nanobox-io/nanobox-golang-stylish"
	
	"github.com/nanobox-io/nanobox/processor"
)

type platformStop struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("platform_stop", platformStopFunc)
}

func platformStopFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return platformStop{control}, nil
}

func (self platformStop) Results() processor.ProcessControl {
	return self.control
}

func (self platformStop) Process() error {

	if err := self.stopServices(); err != nil {
		return err
	}

	return nil
}

// stopServices will stop all the platform services
func (self *platformStop) stopServices() error {
	self.control.Display(stylish.Bullet("Stopping Platform Services"))
	for _, service := range PlatformServices {
		if err := self.stopService(service); err != nil {
			return err
		}
	}

	return nil
}

// stopService will stop an individual service
func (self *platformStop) stopService(service PlatformService) error {

	config := processor.ProcessControl{
		DevMode: self.control.DevMode,
		Verbose: self.control.Verbose,
		DisplayLevel: self.control.DisplayLevel+1,
		Meta: map[string]string{
			"label": service.label,
			"name":  service.name,
		},
	}

	// stop
	if err := processor.Run("service_stop", config); err != nil {
		return err
	}

	return nil
}
