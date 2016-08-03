package platform

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processPlatformDeploy ...
type processPlatformDeploy struct {
	control processor.ProcessControl
}

// this sets up the necessary pieces to do a deploy locally
// which requires a warehouse as well as a portal
func init() {
	processor.Register("platform_deploy", platformDeployFn)
}

//
func platformDeployFn(control processor.ProcessControl) (processor.Processor, error) {
	return processPlatformDeploy{control}, nil
}

//
func (platformDeploy processPlatformDeploy) Results() processor.ProcessControl {
	return platformDeploy.control
}

//
func (platformDeploy processPlatformDeploy) Process() error {

	if err := platformDeploy.provisionServices(); err != nil {
		return err
	}

	return nil
}

// provisionServices will provision all the platform services
func (platformDeploy processPlatformDeploy) provisionServices() error {
	platformDeploy.control.Display(stylish.Bullet("Provisioning Platform Services"))
	for _, service := range DeployServices {
		if err := platformDeploy.provisionService(service); err != nil {
			return err
		}
	}

	return nil
}

// provisionService will provision an individual service
func (platformDeploy processPlatformDeploy) provisionService(service Service) error {

	control := platformDeploy.control.Dup()
	control.DisplayLevel++
	control.Meta["label"] = service.label
	control.Meta["name"] = service.name
	control.Meta["image"] = service.image

	if platformDeploy.isServiceActive(service.name) {
		// start the service if the service is already active
		return processor.Run("service_start", control)
	}

	// otherwise
	// setup the service
	if err := processor.Run("service_setup", control); err != nil {
		return err
	}

	// and configure it
	return processor.Run("service_configure", control)
}

// isServiceActive returns true if a service is already active
func (platformDeploy processPlatformDeploy) isServiceActive(id string) bool {

	// service db entry
	service := models.Service{}

	// fetch the entry from the database, ignoring any errors as the service
	// might not exist yet
	bucket := fmt.Sprintf("%s_%s", config.AppID(), platformDeploy.control.Env)
	data.Get(bucket, id, &service)

	return service.State == "active"
}
