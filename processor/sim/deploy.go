package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processSimDeploy ...
type processSimDeploy struct {
	control processor.ProcessControl
	box     boxfile.Boxfile
}

// deploys the code to the warehouse and builds
// webs, workers, and updates services
// then updates the router with the new code services
func init() {
	processor.Register("sim_deploy", simDeployFn)
}

//
func simDeployFn(control processor.ProcessControl) (processor.Processor, error) {
	// control.Meta["processSimDeploy-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return processSimDeploy{control: control}, nil
}

//
func (simDeploy processSimDeploy) Results() processor.ProcessControl {
	return simDeploy.control
}

//
func (simDeploy processSimDeploy) Process() error {
	// set the mode of this processes
	// this allows the dev and the deploy to be isolated
	simDeploy.control.Env = "sim"

  // run the share init which gives access to docker
  if err := processor.Run("env_init", simDeploy.control); err != nil {
    return err
  }

	// get the platform deploy ready
	if err := processor.Run("platform_deploy", simDeploy.control); err != nil {
		return err
	}

	if err := simDeploy.publishCode(); err != nil {
		return err
	}

	// remove all the previous code services
	if err := processor.Run("code_clean", simDeploy.control); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	if err := processor.Run("service_sync", simDeploy.control); err != nil {
		return err
	}

	// start code
	if err := simDeploy.startCodeServices(); err != nil {
		return err
	}

	// update nanoagent portal
	if err := processor.Run("update_portal", simDeploy.control); err != nil {
		return err
	}

	// complete message
	fmt.Println("The deploy completed successfully!")

	return nil
}

// publishCode ...
func (simDeploy *processSimDeploy) publishCode() error {

	// setup the var's required for code_publish
	hoarder := models.Service{}
	bucket := fmt.Sprintf("%s_%s", config.AppName(), simDeploy.control.Env)
	data.Get(bucket, "hoarder", &hoarder)

	simDeploy.control.Meta["build_id"] = "1234"
	simDeploy.control.Meta["warehouse_url"] = hoarder.InternalIP
	simDeploy.control.Meta["warehouse_token"] = "123"

	// publish code
	publishProcessor, err := processor.Build("code_publish", simDeploy.control)
	if err != nil {
		return err
	}

	err = publishProcessor.Process()
	if err != nil {
		return err
	}

	publishResult := publishProcessor.Results()
	if publishResult.Meta["boxfile"] == "" {
		return fmt.Errorf("publishCode: the boxfile was empty")
	}

	// store the boxfile on mydeploy
	simDeploy.box = boxfile.New([]byte(publishResult.Meta["boxfile"]))

	// set it in the control file so child process have access to it as well
	simDeploy.control.Meta["boxfile"] = publishResult.Meta["boxfile"]

	return nil
}

// startCodeServices ...
func (simDeploy *processSimDeploy) startCodeServices() error {
	code := processor.ProcessControl{
		Env: simDeploy.control.Env,
		Verbose: simDeploy.control.Verbose,
		Meta: map[string]string{
			"boxfile":         simDeploy.control.Meta["boxfile"],
			"build_id":        simDeploy.control.Meta["build_id"],
			"warehouse_url":   simDeploy.control.Meta["warehouse_url"],
			"warehouse_token": simDeploy.control.Meta["warehouse_token"],
		},
	}

	// synchronize my code services
	if err := processor.Run("code_sync", code); err != nil {
		return err
	}

	return nil
}
