package platform

import (
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processPlatformSetup ...
type processPlatformSetup struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("platform_setup", platformSetupFn)
}

//
func platformSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return processPlatformSetup{control}, nil
}

//
func (platformSetup processPlatformSetup) Results() processor.ProcessControl {
	return platformSetup.control
}

//
func (platformSetup processPlatformSetup) Process() error {

	if err := platformSetup.provisionServices(); err != nil {
		return err
	}

	return nil
}

// provisionServices will provision all the platform services
func (platformSetup processPlatformSetup) provisionServices() error {
	platformSetup.control.Display(stylish.Bullet("Provisioning Platform Services"))
	for _, service := range Services {
		if err := platformSetup.provisionService(service); err != nil {
			return err
		}
	}

	return nil
}

// provisionService will provision an individual service
func (platformSetup processPlatformSetup) provisionService(service Service) error {

	config := processor.ProcessControl{
		DevMode:      platformSetup.control.DevMode,
		Verbose:      platformSetup.control.Verbose,
		DisplayLevel: platformSetup.control.DisplayLevel + 1,
		Meta: map[string]string{
			"label": service.label,
			"name":  service.name,
			"image": service.image,
		},
	}

	if platformSetup.isServiceActive(service.name) {
		// start the service if the service is already active
		return processor.Run("service_start", config)
	}

	// otherwise
	// setup the service
	if err := processor.Run("service_setup", config); err != nil {
		return err
	}

	// and configure it
	return processor.Run("service_configure", config)
}

// isServiceActive returns true if a service is already active
func (platformSetup processPlatformSetup) isServiceActive(id string) bool {

	// service db entry
	service := models.Service{}

	// fetch the entry from the database, ignoring any errors as the service
	// might not exist yet
	data.Get(config.AppName(), id, &service)

	return service.State == "active"
}
