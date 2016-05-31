package processor

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type devDeploy struct {
	control ProcessControl
	box    boxfile.Boxfile
}

func init() {
	Register("dev_deploy", devDeployFunc)
}

func devDeployFunc(control ProcessControl) (Processor, error) {
	// control.Meta["devDeploy-control"]
	// do some control validation
	// check on the meta for the flags and make sure they work

	return devDeploy{control: control}, nil
}

func (self devDeploy) Results() ProcessControl {
	return self.control
}

func (self devDeploy) Process() error {

	// defer the clean up so if we exit early the
	// cleanup will always happen
	defer func() {
		if err := Run("dev_teardown", self.control); err != nil {
			// this is bad, really bad...
			// we should probably print a pretty message explaining that the app
			// was left in a bad state and needs to be reset
			os.Exit(1)
		}
	}()

	if err := Run("dev_setup", self.control); err != nil {
		// todo: how to display this?
		return err
	}

	if err := self.publishCode(); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	if err := Run("service_sync", self.control); err != nil {
		return err
	}

	// start code
	if err := self.startCodeServices(); err != nil {
		return err
	}

	// clean up the code services
	defer func() {
		if err := Run("code_clean", self.control); err != nil {
			// output this error message
			// it doesnt break anything if the clean fails.
		}
	}()

	// update nanoagent portal
	if err := Run("update_portal", self.control); err != nil {
		return err
	}

	// hang and do some logging until they are done
	if err := Run("mist_log", self.control); err != nil {
		return err
	}

	return nil
}

func (self *devDeploy) publishCode() error {

	// setup the var's required for code_publish
	hoarder := models.Service{}
	data.Get(util.AppName(), "hoarder", &hoarder)
	self.control.Meta["build_id"] = "1234"
	self.control.Meta["warehouse_url"] = hoarder.InternalIP
	self.control.Meta["warehouse_token"] = "123"

	// publish code
	publishProcessor, err := Build("code_publish", self.control)
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
	// store the boxfile on myself
	self.box = boxfile.New([]byte(publishResult.Meta["boxfile"]))

	// set it in the control file so child process have access to it as well
	self.control.Meta["boxfile"] = publishResult.Meta["boxfile"]
	return nil
}

func (self *devDeploy) startCodeServices() error {
	code := ProcessControl{
		DevMode: self.control.DevMode,
		Verbose: self.control.Verbose,
		Meta: map[string]string{
			"boxfile":         self.control.Meta["boxfile"],
			"build_id":        self.control.Meta["build_id"],
			"warehouse_url":   self.control.Meta["warehouse_url"],
			"warehouse_token": self.control.Meta["warehouse_token"],
		},
	}

	// synchronize my code services
	if err := Run("code_sync", code); err != nil {
		return err
	}
	return nil
}
