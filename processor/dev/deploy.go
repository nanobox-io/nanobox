package dev

import (
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/processor"
)

// processDevDeploy ...
type processDevDeploy struct {
	control processor.ProcessControl
	box     boxfile.Boxfile
}

//
func init() {
	processor.Register("dev_deploy", devDeployFn)
}

//
func devDeployFn(control processor.ProcessControl) (processor.Processor, error) {
	// control.Meta["processDevDeploy-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return processDevDeploy{control: control}, nil
}

//
func (devDeploy processDevDeploy) Results() processor.ProcessControl {
	return devDeploy.control
}

//
func (devDeploy processDevDeploy) Process() error {
	// set the mode of this processes
	// this allows the dev and the deploy to be isolated
	devDeploy.control.Env = "dev"

	// run the share init which gives access to docker
	if err := processor.Run("env_init", devDeploy.control); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	if err := processor.Run("service_sync", devDeploy.control); err != nil {
		return err
	}

	// complete message
	devDeploy.control.Display(stylish.Bullet("Deploy complete!"))

	return nil
}
