package processors

import (
	"fmt"
	
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/util/odin"
)

//
func Deploy(envModel *models.Env, deployConfig DeployConfig) error {

	// commented until we have the --force as a flag
	// if envModel.DeployedID != "" && envModel.BuiltID == envModel.DeployedID {
	// 	// shortcut if we have already deployed
	// 	return nil
	// }

	// find the app id
	appID := models.AppIDByAlias(deployConfig.App)
	if appID == "" {
		// todo: better messaging informing that we couldn't find a link
		return fmt.Errorf("app is not linked")
	}

	warehouseConfig, err := getWarehouseConfig(envModel, appID)
	if err != nil {
		return fmt.Errorf("unable to generate warehouse config: %s", err.Error())
	}

	// publish to remote warehouse
	if err := code.Publish(envModel, warehouseConfig); err != nil {
		return fmt.Errorf("failed to publish build to app's warehouse: %s", err.Error())
	}

	// tell odin what happened
	if err := odin.Deploy(appID, warehouseConfig.BuildID, envModel.BuiltBoxfile, deployConfig.Message); err != nil {
		lumber.Error("deploy:odin.Deploy(%s,%s,%s,%s): %s", appID, warehouseConfig.BuildID, envModel.BuiltBoxfile, deployConfig.Message, err.Error())
		return fmt.Errorf("failed to deploy code to app: %s", err.Error())
	}

	envModel.DeployedID = envModel.BuiltID
	if err := envModel.Save(); err != nil {
		lumber.Error("deploy:models:Env:Save(): %s", err.Error())
		return fmt.Errorf("failed to save build ID: %s", err.Error())
	}

	return nil
}

// setWarehouseToken ...
func getWarehouseConfig(envModel *models.Env, appID string) (warehouseConfig code.WarehouseConfig, err error) {

	token, url, err := odin.GetWarehouse(appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetWarehouse(%s): %s", appID, err.Error())
		err = fmt.Errorf("failed to fetch warehouse information from nanobox: %s", err.Error())
		return
	}

	// get the previous build if there was one
	prevBuild, err := odin.GetPreviousBuild(appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetPreviousBuild(%s): %s", appID, err.Error())
		err = fmt.Errorf("failed to query previous deploys from nanobox: %s", err.Error())
		return
	}
	
	warehouseConfig.BuildID        = envModel.BuiltID
	warehouseConfig.WarehouseURL   = url
	warehouseConfig.WarehouseToken = token
	warehouseConfig.PreviousBuild  = prevBuild
	
	return
}
