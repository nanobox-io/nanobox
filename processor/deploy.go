package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"
)

// processDeploy ...
type processDeploy struct {
	control ProcessControl
}

//
func init() {
	Register("deploy", deployFn)
}

//
func deployFn(control ProcessControl) (Processor, error) {
	// control.Meta["deploy-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return &processDeploy{control}, nil
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

// setWarehouseToken ...
func (deploy *processDeploy) setWarehouseToken() error {
	// get remote hoarder credentials
	deploy.control.Meta["app_id"] = getAppID(deploy.control.Meta["alias"])
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
