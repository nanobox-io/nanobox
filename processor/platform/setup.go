package platform

import (
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/processor"
)

type platformSetup struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("platform_setup", platformSetupFunc)
}

func platformSetupFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return platformSetup{control}, nil
}

func (self platformSetup) Results() processor.ProcessControl {
	return self.control
}

func (self platformSetup) Process() error {

	if err := self.provisionServices(); err != nil {
		return err
	}

	return nil
}

// provisionServices will provision all the platform services
func (self platformSetup) provisionServices() error {
	self.control.Display(stylish.Bullet("Provisioning Platform Services"))
	for _, service := range PlatformServices {
		if err := self.provisionService(service); err != nil {
			return err
		}
	}

	return nil
}

// provisionService will provision an individual service
func (self platformSetup) provisionService(service PlatformService) error {

	config := processor.ProcessControl{
		DevMode:      self.control.DevMode,
		Verbose:      self.control.Verbose,
		DisplayLevel: self.control.DisplayLevel + 1,
		Meta: map[string]string{
			"label": service.label,
			"name":  service.name,
			"image": service.image,
		},
	}

	if self.isServiceActive(service.name) {
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
func (self platformSetup) isServiceActive(id string) bool {

	// service db entry
	service := models.Service{}

	// fetch the entry from the database, ignoring any errors as the service
	// might not exist yet
	data.Get(util.AppName(), id, &service)

	return service.State == "active"
}
