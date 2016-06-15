package dev

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDevDeploy ...
type processDevDeploy struct {
	control processor.ProcessControl
	box     boxfile.Boxfile
}

//
func init() {
	processor.Register("dev_deploy", devDeployFunc)
}

//
func devDeployFunc(control processor.ProcessControl) (processor.Processor, error) {
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

	// defer the clean up so if we exit early the
	// cleanup will always happen
	defer func() {
		if err := processor.Run("dev_teardown", devDeploy.control); err != nil {
			// this is bad, really bad...
			// we should probably print a pretty message explaining that the app
			// was left in a bad state and needs to be reset
			os.Exit(1)
		}
	}()

	if err := processor.Run("dev_setup", devDeploy.control); err != nil {
		// todo: how to display this?
		return err
	}

	if err := devDeploy.publishCode(); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	if err := processor.Run("service_sync", devDeploy.control); err != nil {
		return err
	}

	// start code
	if err := devDeploy.startCodeServices(); err != nil {
		return err
	}

	// clean up the code services
	defer func() {
		if err := processor.Run("code_clean", devDeploy.control); err != nil {
			// output this error message
			// it doesnt break anything if the clean fails.
		}
	}()

	// update nanoagent portal
	if err := processor.Run("update_portal", devDeploy.control); err != nil {
		return err
	}

	// hang and do some logging until they are done
	if err := processor.Run("mist_log", devDeploy.control); err != nil {
		return err
	}

	return nil
}

// publishCode ...
func (devDeploy *processDevDeploy) publishCode() error {

	// setup the var's required for code_publish
	hoarder := models.Service{}
	data.Get(config.AppName(), "hoarder", &hoarder)
	devDeploy.control.Meta["build_id"] = "1234"
	devDeploy.control.Meta["warehouse_url"] = hoarder.InternalIP
	devDeploy.control.Meta["warehouse_token"] = "123"

	// publish code
	publishProcessor, err := processor.Build("code_publish", devDeploy.control)
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
	devDeploy.box = boxfile.New([]byte(publishResult.Meta["boxfile"]))

	// set it in the control file so child process have access to it as well
	devDeploy.control.Meta["boxfile"] = publishResult.Meta["boxfile"]

	return nil
}

// startCodeServices ...
func (devDeploy *processDevDeploy) startCodeServices() error {
	code := processor.ProcessControl{
		DevMode: devDeploy.control.DevMode,
		Verbose: devDeploy.control.Verbose,
		Meta: map[string]string{
			"boxfile":         devDeploy.control.Meta["boxfile"],
			"build_id":        devDeploy.control.Meta["build_id"],
			"warehouse_url":   devDeploy.control.Meta["warehouse_url"],
			"warehouse_token": devDeploy.control.Meta["warehouse_token"],
		},
	}

	// synchronize my code services
	if err := processor.Run("code_sync", code); err != nil {
		return err
	}

	return nil
}
