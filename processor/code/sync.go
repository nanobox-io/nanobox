package code

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processCodeSync is used by the deploy process to syncronize the code parts
// from the boxfile
type processCodeSync struct {
	control processor.ProcessControl
	box     boxfile.Boxfile
}

//
func init() {
	processor.Register("code_sync", codeSyncFn)
}

//
func codeSyncFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// {"build":"%s","warehouse":"%s","warehouse_token":"123","boxfile":"%s"}
	if control.Meta["build_id"] == "" ||
		control.Meta["boxfile"] == "" ||
		control.Meta["warehouse_url"] == "" ||
		control.Meta["warehouse_token"] == "" {
		return nil, fmt.Errorf("missing boxfile || build_id || warehouse_url || warehouse_token")
	}

	return &processCodeSync{control: control}, nil
}

//
func (codeSync processCodeSync) Results() processor.ProcessControl {
	return codeSync.control
}

//
func (codeSync *processCodeSync) Process() error {
	// increment the counter so we know how many deploys are waiting
	counter.Increment(config.AppName() + "_deploy")

	// do not allow more then one process to run the code sync at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// set the boxfile and make sure its valid
	if err := codeSync.setBoxfile(); err != nil {
		return err
	}
	fmt.Printf("%+v\n", codeSync.box)

	// iterate over the code nodes and build containers for each of them
	for _, codeName := range codeSync.box.Nodes("code") {
		// pull the image from the boxfile; default to a reasonable alternative if
		// none is given
		image := codeSync.box.Node(codeName).StringValue("image")
		if image == "" {
			image = "nanobox/code:v1"
		}

		// create a new process config for code ensuring it has access to the warehouse
		// and the boxfile
		code := processor.ProcessControl{
			Env: codeSync.control.Env,
			Verbose: codeSync.control.Verbose,
			Meta: map[string]string{
				"name":            codeName,
				"image":           image,
				"boxfile":         codeSync.control.Meta["boxfile"],
				"build_id":        codeSync.control.Meta["build_id"],
				"warehouse_url":   codeSync.control.Meta["warehouse_url"],
				"warehouse_token": codeSync.control.Meta["warehouse_token"],
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
			return fmt.Errorf("code_configure (%s): %s\n", codeName, err.Error())
		}
	}

	return nil
}

// setBoxfile ...
func (codeSync *processCodeSync) setBoxfile() error {
	codeSync.box = boxfile.New([]byte(codeSync.control.Meta["boxfile"]))
	if !codeSync.box.Valid {
		return fmt.Errorf("Invalid Boxfile")
	}

	return nil
}
