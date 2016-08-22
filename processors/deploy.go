package processors

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Deploy ...
type Deploy struct {
	Env     models.Env
	App     string
	Message string

	appID          string
	buildID        string
	warehouseURL   string
	warehouseToken string
	previousBuild  string
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
	if err := odin.Deploy(deploy.appID, deploy.buildID, deploy.Env.BuiltBoxfile, deploy.Message); err != nil {
		lumber.Error("deploy:odin.Deploy(%s,%s,%s,%s): %s", deploy.appID, deploy.buildID, deploy.Env.BuiltBoxfile, deploy.Message, err.Error())
		return err
	}
	return nil
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
		lumber.Error("deploy:setWarehouseToken:GetWarehouse(%s): %s", deploy.appID, err.Error())
		return
	}

	// get the previous build if there was one
	deploy.previousBuild, err = odin.GetPreviousBuild(deploy.appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetPreviousBuild(%s): %s", deploy.appID, err.Error())
		return
	}
	return
}

// publishCode ...
func (deploy *Deploy) publishCode() error {

	codePublish := code.Publish{
		Env:            deploy.Env,
		BuildID:        deploy.buildID,
		WarehouseURL:   deploy.warehouseURL,
		WarehouseToken: deploy.warehouseToken,
		PreviousBuild:  deploy.previousBuild,
	}

	return codePublish.Run()
}
