package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/processor/code"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/processor/platform"
)

// deploys the code to the warehouse and builds
// webs, workers, and updates services
// then updates the router with the new code services
type Deploy struct {
	// mandatory
	App models.App
	// added
	Env models.Env
}

//
func (deploy Deploy) Run() error {

	deploy.Env, _ = models.FindEnvByID(deploy.App.EnvID)

	// run the share init which gives access to docker
	envInit := env.Init{}
	if err := envInit.Run(); err != nil {
		return err
	}

	platformSetup := platform.Setup{App: deploy.App}
	if err := platformSetup.Run(); err != nil {
		return err
	}

	if err := deploy.publishCode(); err != nil {
		return err
	}

	codeClean := code.Clean{App: deploy.App}
	// remove all the previous code services
	if err := codeClean.Run(); err != nil {
		return err
	}

	componentSync := component.Sync{
		Env: deploy.Env,
		App: deploy.App,
	}
	// syncronize the services as per the new boxfile
	if err := componentSync.Run(); err != nil {
		return err
	}

	// start code
	if err := deploy.startCodeServices(); err != nil {
		return err
	}

	if err := deploy.runDeployHook("before_deploy"); err != nil {
		return err
	}

	// update nanoagent portal
	platformUpdatePortal := platform.UpdatePortal{App: deploy.App}
	if err := platformUpdatePortal.Run(); err != nil {
		return err
	}

	if err := deploy.runDeployHook("after_deploy"); err != nil {
		return err
	}

	// complete message
	fmt.Println("The deploy completed successfully!")

	return nil
}

// publishCode ...
func (deploy *Deploy) publishCode() error {

	// setup the var's required for code_publish
	hoarder, _ := models.FindComponentBySlug(deploy.App.ID, "hoarder")

	codePublish := code.Publish{
		Env: deploy.Env,
		BuildID: "1234",
		WarehouseURL: hoarder.InternalIP,
		WarehouseToken: "123",
	}

	return codePublish.Run()
}

// startCodeServices ...
func (deploy *Deploy) startCodeServices() error {

	// synchronize my code services
	hoarder, _ := models.FindComponentBySlug(deploy.App.ID, "hoarder")

	codeSync := code.Sync{
		App: deploy.App,
		BuildID: "1234",
		WarehouseURL: hoarder.InternalIP,
		WarehouseToken: "123",
	}

	return codeSync.Run()
}

// run the before/after hooks and populate the necessary data
func (deploy *Deploy) runDeployHook(hookType string) error {
	box := boxfile.New([]byte(deploy.App.DeployedBoxfile))

	// run the hooks for each service in the boxfile
	for _, componentName := range box.Nodes("code") {

		component, err := models.FindComponentBySlug(deploy.App.ID, componentName)
		if err != nil {
			// no component for that thing in the database..
			// prolly need to report this error but we might not want to fail
			continue 
		}

		deployHook := DeployHook{
			App: deploy.App,
			Component: component,
			HookType: hookType,
		}
		if err := deployHook.Run(); err != nil {
			return err
		}
	}

	return nil
}
