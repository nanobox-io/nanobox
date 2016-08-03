package platform

import (
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
)

//
type processPlatformStop struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("platform_stop", platformStopFn)
}

//
func platformStopFn(control processor.ProcessControl) (processor.Processor, error) {
	return processPlatformStop{control}, nil
}

//
func (platformStop processPlatformStop) Results() processor.ProcessControl {
	return platformStop.control
}

//
func (platformStop processPlatformStop) Process() error {
	return platformStop.stopServices()
}

// stopServices will stop all the platform services
func (platformStop *processPlatformStop) stopServices() error {

	// stop all the platform services weather they have been created
	// or not. (some may not be created in all instances)
	platformStop.control.Display(stylish.Bullet("Stopping Platform Services..."))
	for _, service := range append(SetupServices, DeployServices...) {
		if err := platformStop.stopService(service); err != nil {
			return err
		}
	}

	return nil
}

// stopService will stop an individual service
func (platformStop *processPlatformStop) stopService(service Service) error {
	//
	control := platformStop.control.Dup()
	control.DisplayLevel++
	control.Meta["label"] = service.label
	control.Meta["name"] = service.name

	//
	return processor.Run("service_stop", control)
}
