package code

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
)

// used by the deploy process to syncronize the code parts from the boxfile
type codeSync struct {
	control processor.ProcessControl
	box    boxfile.Boxfile
}

func init() {
	processor.Register("code_sync", codeSyncFunc)
}

func codeSyncFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// {"build":"%s","warehouse":"%s","warehouse_token":"123","boxfile":"%s"}
	if control.Meta["build_id"] == "" ||
		control.Meta["boxfile"] == "" ||
		control.Meta["warehouse_url"] == "" ||
		control.Meta["warehouse_token"] == "" {
		return nil, fmt.Errorf("missing boxfile || build_id || warehouse_url || warehouse_token")
	}
	return &codeSync{control: control}, nil
}

func (self codeSync) Results() processor.ProcessControl {
	return self.control
}

func (self *codeSync) Process() error {
	// increment the counter so we know how many deploys are waiting
	counter.Increment(util.AppName() + "_deploy")

	// do not allow more then one process to run the code sync at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// set the boxfile and make sure its valid
	if err := self.setBoxfile(); err != nil {
		return err
	}

	// iterate over the code nodes and build containers for each of them
	for _, codeName := range self.box.Nodes("code") {
		// pull the image from the boxfile.
		// default to a reasonable alternative if non is given
		image := self.box.Node(codeName).StringValue("image")
		if image == "" {
			image = "nanobox/code"
		}

		// create a new process config for code
		// ensuring it has access to the warehouse
		// and the boxfile
		code := processor.ProcessControl{
			DevMode: self.control.DevMode,
			Verbose: self.control.Verbose,
			Meta: map[string]string{
				"name":            codeName,
				"image":           image,
				"boxfile":         self.control.Meta["boxfile"],
				"build_id":        self.control.Meta["build_id"],
				"warehouse_url":   self.control.Meta["warehouse_url"],
				"warehouse_token": self.control.Meta["warehouse_token"],
			},
		}

		// run the code setup process with the new config
		err := processor.Run("code_setup", code)
		if err != nil {
			return fmt.Errorf("code_setup (%s): %s\n", codeName, err.Error())
		}

		// configure this code container
		err = processor.Run("code_configure", code)
		if err != nil {
			return fmt.Errorf("code_start (%s): %s\n", codeName, err.Error())
		}

	}
	return nil
}

func (self *codeSync) setBoxfile() error {
	self.box = boxfile.New([]byte(self.control.Meta["boxfile"]))
	if !self.box.Valid {
		return fmt.Errorf("Invalid Boxfile")
	}
	return nil
}
