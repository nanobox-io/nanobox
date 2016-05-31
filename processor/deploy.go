package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/production_api"
)

type deploy struct {
	control ProcessControl
}

func init() {
	Register("deploy", deployFunc)
}

func deployFunc(control ProcessControl) (Processor, error) {
	// control.Meta["deploy-control"]
	// do some control validation
	// check on the meta for the flags and make sure they work

	return &deploy{control}, nil
}

func (self deploy) Results() ProcessControl {
	return self.control
}

func (self *deploy) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.control)
	if err != nil {
		return err
	}


	if err := self.setWarehouseToken(); err != nil {
		return err
	}

	if err := self.publishCode(); err != nil {
		return err
	}

	// tell odin what happened
	return production_api.Deploy(self.control.Meta["app_id"], self.control.Meta["build_id"], self.control.Meta["boxfile"], self.control.Meta["message"])
}

func (self *deploy) setWarehouseToken() error {
	// get remote hoarder credentials
	self.control.Meta["app_id"] = getAppID(self.control.Meta["alias"])
	self.control.Meta["build_id"] = util.RandomString(30)
	warehouseToken, warehouseUrl, err := production_api.GetWarehouse(self.control.Meta["app_id"])
	if err != nil {
		return err
	}
	self.control.Meta["warehouse_token"] = warehouseToken
	self.control.Meta["warehouse_url"] = warehouseUrl
	return nil
}

func (self *deploy) publishCode() error {
	publishProcessor, err := Build("code_publish", self.control)
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
	self.control.Meta["boxfile"] = publishResult.Meta["boxfile"]
	return nil
}
