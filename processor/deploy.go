package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"
)

// processDeploy ...
type processDeploy struct {
	control ProcessControl
	app     string
}

//
func init() {
	Register("deploy", deployFn)
}

//
func deployFn(control ProcessControl) (Processor, error) {
	deploy := &processDeploy{control: control}
	return deploy, deploy.validateMeta()
}

//
func (deploy processDeploy) Results() ProcessControl {
	return deploy.control
}

//
func (deploy *processDeploy) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", deploy.control)
	if err != nil {
		return err
	}

	if err := deploy.setWarehouseToken(); err != nil {
		return err
	}

	if err := deploy.publishCode(); err != nil {
		return err
	}

	// tell odin what happened
	return odin.Deploy(deploy.control.Meta["app_id"], deploy.control.Meta["build_id"], deploy.control.Meta["boxfile"], deploy.control.Meta["message"])
}

// validateMeta validates that the required metadata exists
func (deploy *processDeploy) validateMeta() error {

	// set app (required) and ensure it's provided
	deploy.app = deploy.control.Meta["app"]
	if deploy.app == "" {
		return fmt.Errorf("Missing required meta value 'app'")
	}

	return nil
}

// setWarehouseToken ...
func (deploy *processDeploy) setWarehouseToken() error {

	// get remote hoarder credentials
	deploy.control.Meta["app_id"] = getAppID(deploy.app)
	deploy.control.Meta["build_id"] = util.RandomString(30)
	warehouseToken, warehouseURL, err := odin.GetWarehouse(deploy.control.Meta["app_id"])
	if err != nil {
		return err
	}
	
	deploy.control.Meta["warehouse_token"] = warehouseToken
	deploy.control.Meta["warehouse_url"] = warehouseURL
	return nil
}

// publishCode ...
func (deploy *processDeploy) publishCode() error {
	publishProcessor, err := Build("code_publish", deploy.control)
	if err != nil {
		return err
	}

	if err := publishProcessor.Process(); err != nil {
		return err
	}
	publishResult := publishProcessor.Results()
	if publishResult.Meta["boxfile"] == "" {
		return fmt.Errorf("the boxfile from publish was blank")
	}
	// boxfile := boxfile.New([]byte(publishResult.Meta["boxfile"]))
	deploy.control.Meta["boxfile"] = publishResult.Meta["boxfile"]

	return nil
}
