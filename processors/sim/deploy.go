package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/platform"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
)

// deploys the code to the warehouse and builds
// webs, workers, and updates services
// then updates the router with the new code services
func Deploy(envModel *models.Env, appModel *models.App) error {
	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	if err := platform.Deploy(appModel); err != nil {
		return err
	}

	// create the warehouse config for child processes
	hoarder, _ := models.FindComponentBySlug(appModel.ID, "hoarder")

	warehouseConfig := code.WarehouseConfig{
		BuildID:        "1234",
		WarehouseURL:   hoarder.InternalIP,
		WarehouseToken: "123",
	}

	// publish the code
	if err := code.Publish(envModel, warehouseConfig); err != nil {
		return fmt.Errorf("unable to publish code: %s", err.Error())
	}

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return fmt.Errorf("unable to synchronize data components: %s", err.Error())
	}

	// start code
	if err := code.Sync(appModel, warehouseConfig); err != nil {
		return fmt.Errorf("failed to add code components: %s", err.Error())
	}

	if err := finalizeDeploy(appModel); err != nil {
		return fmt.Errorf("failed to finalize deploy: %s", err.Error())
	}

	return nil
}

// update the router and run deploy hooks
func finalizeDeploy(appModel *models.App) error {
	display.OpenContext("Finalizing deploy")
	defer display.CloseContext()

	display.StartTask("Running before_deploy hooks")
	if err := runDeployHook(appModel, "before_deploy"); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run before deploy hooks: %s", err.Error())
	}
	display.StopTask()

	// update nanoagent portal
	display.StartTask("Updating router")
	if err := platform.UpdatePortal(appModel); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to update router: %s", err.Error())
	}
	display.StopTask()

	display.StartTask("Running after_deploy hooks")
	if err := runDeployHook(appModel, "after_deploy"); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run after deloy hooks: %s", err.Error())
	}
	display.StopTask()

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
