package processors

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

//
func Deploy(envModel *models.Env, deployConfig DeployConfig) error {

	appID := deployConfig.App

	// fetch the link
	link, ok := envModel.Links[deployConfig.App]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(link.Endpoint)
		// set the app id
		appID = link.ID
	}

	// if an endpoint was provided as a flag, override the linked endpoint
	if deployConfig.Endpoint != "" {
		odin.SetEndpoint(deployConfig.Endpoint)
	}

	// set the app id to the directory name if it's default
	if appID == "default" {
		appID = config.AppName()
	}

	// validate access to the app
	if err := helpers.ValidateOdinApp(appID); err != nil {
		// the validation already printed the error
		return nil
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

	fmt.Printf("%s Deploy was successufully submitted! Check your dashboard for progress.\n", display.TaskComplete)

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

	warehouseConfig.BuildID = envModel.BuiltID
	warehouseConfig.WarehouseURL = url
	warehouseConfig.WarehouseToken = token
	warehouseConfig.PreviousBuild = prevBuild

	return
}
