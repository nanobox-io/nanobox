package processors

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"
)

//
func Deploy(envModel *models.Env, deployConfig DeployConfig) error {

	// setup the environment (boot vm)
	if err := provider.Setup(); err != nil {
		return err
	}

	appID := getAppID(deployConfig.App)

	warehouseConfig, err := getWarhouseConfig(appID)
	if err != nil {
		return err
	}

	// publish to remote warehouse
	if err := code.Publish(envModel, warehouseConfig); err != nil {
		return err
	}

	// tell odin what happened
	if err := odin.Deploy(appID, warehouseConfig.BuildID, envModel.BuiltBoxfile, deployConfig.Message); err != nil {
		lumber.Error("deploy:odin.Deploy(%s,%s,%s,%s): %s", appID, warehouseConfig.BuildID, envModel.BuiltBoxfile, deployConfig.Message, err.Error())
		return err
	}
	return nil
}

// setWarehouseToken ...
func getWarhouseConfig(appID string) (warehouseConfig code.WarehouseConfig, err error) {

	token, url, err := odin.GetWarehouse(appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetWarehouse(%s): %s", appID, err.Error())
		return
	}

	// get the previous build if there was one
	prevBuild, err := odin.GetPreviousBuild(appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetPreviousBuild(%s): %s", appID, err.Error())
		return
	}
	warehouseConfig.BuildID = util.RandomString(30)
	warehouseConfig.WarehouseURL = url
	warehouseConfig.WarehouseToken = token
	warehouseConfig.PreviousBuild = prevBuild
	return
}
