package processors

import (
	"fmt"
	
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/odin"
)

//
func Deploy(envModel *models.Env, deployConfig DeployConfig) error {

	if envModel.DeployedID != "" && envModel.BuiltID == envModel.DeployedID {
		// shortcut if we have already deployed
		return nil
	}
	// setup the environment (boot vm)
	if err := provider.Setup(); err != nil {
		return err
	}

	appID := getAppID(deployConfig.App)

	warehouseConfig, err := getWarhouseConfig(envModel, appID)
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

	envModel.DeployedID = envModel.BuiltID
	if err := envModel.Save(); err != nil {
		lumber.Error("deploy:models:Env:Save(): %s", err.Error())
		return fmt.Errorf("env model: %s", err.Error())
	}

	return nil
}

// setWarehouseToken ...
func getWarhouseConfig(envModel *models.Env, appID string) (warehouseConfig code.WarehouseConfig, err error) {

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
	warehouseConfig.BuildID = envModel.BuiltID
	warehouseConfig.WarehouseURL = url
	warehouseConfig.WarehouseToken = token
	warehouseConfig.PreviousBuild = prevBuild
	return
}
