package sim

import (
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/platform"
	"github.com/nanobox-io/nanobox/util/display"
)

// deploys the code to the warehouse and builds
// webs, workers, and updates services
// then updates the router with the new code services
func Deploy(envModel *models.Env, appModel *models.App) error {
	display.OpenContext("Deploying Sim")
	defer display.CloseContext()

	// run the share init which gives access to docker
	if err := provider.Init(); err != nil {
		return err
	}

	display.StartTask("starting services for deploy")
	if err := platform.Deploy(appModel); err != nil {
		return err
	}
	display.StopTask()

	// create the warehouse config for child processes
	hoarder, _ := models.FindComponentBySlug(appModel.ID, "hoarder")

	warehouseConfig := code.WarehouseConfig{
		BuildID:        "1234",
		WarehouseURL:   hoarder.InternalIP,
		WarehouseToken: "123",
	}

	// publish the code
	if err := code.Publish(envModel, warehouseConfig); err != nil {
		return err
	}

	// remove all the previous code services
	if err := code.Clean(appModel); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return err
	}

	// start code
	if err := code.Sync(AppModel, warehouseConfig); err != nil {
		return err
	}

	if err := runDeployHook(appModel, "before_deploy"); err != nil {
		return err
	}

	// update nanoagent portal
	if err := platform.UpdatePortal(appModel); err != nil {
		return err
	}

	if err := runDeployHook(appModel, "after_deploy"); err != nil {
		return err
	}

	// complete message

	return nil
}

// run the before/after hooks and populate the necessary data
func runDeployHook(appModel *models.App, hookType string) error {
	box := boxfile.New([]byte(appModel.DeployedBoxfile))

	// run the hooks for each service in the boxfile
	for _, componentName := range box.Nodes("code") {

		component, err := models.FindComponentBySlug(appModel.ID, componentName)
		if err != nil {
			// no component for that thing in the database..
			// prolly need to report this error but we might not want to fail
			continue
		}

		if err := DeployHook(appModel, component, hookType); err != nil {
			return err
		}
	}

	return nil
}
