package processor

import (

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"
	"github.com/nanobox-io/nanobox/processor/provider"
	"github.com/nanobox-io/nanobox/processor/code"
)

// Deploy ...
type Deploy struct {
	Env models.Env
	App     string
	Message string

	appID   string
	buildID string
	warehouseURL string
	warehouseToken string
	previousBuild string
}

//
func (deploy *Deploy) Run() error {
	// setup the environment (boot vm)
	providerSetup := provider.Setup{}
	if err := providerSetup.Run(); err != nil {
		return err
	}

	if err := deploy.setWarehouseToken(); err != nil {
		return err
	}

	if err := deploy.publishCode(); err != nil {
		return err
	}

	// tell odin what happened
	return odin.Deploy(deploy.appID, deploy.buildID, deploy.Env.BuiltBoxfile, deploy.Message)
}

// setWarehouseToken ...
func (deploy *Deploy) setWarehouseToken() (err error) {

	// get remote hoarder credentials
	deploy.appID = getAppID(deploy.App)
	// TODO: could make this not as random but based on something
	// so if the same code was 'deployed' odin could react??
	deploy.buildID = util.RandomString(30)
	deploy.warehouseToken, deploy.warehouseURL, err = odin.GetWarehouse(deploy.appID)
	if err != nil {
		return
	}

	// get the previous build if there was one
	deploy.previousBuild, err = odin.GetPreviousBuild(deploy.appID)
	return
}

// publishCode ...
func (deploy *Deploy) publishCode() error {

	codePublish := code.Publish{
		Env: deploy.Env,
		BuildID: deploy.buildID,
		WarehouseURL: deploy.warehouseURL,
		WarehouseToken: deploy.warehouseToken,
		PreviousBuild: deploy.previousBuild,
	}

	return codePublish.Run()
}

