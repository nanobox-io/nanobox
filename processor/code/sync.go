package code

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/models"
)

// Sync is used by the deploy process to syncronize the code parts
// from the boxfile
type Sync struct {
	App models.App
	Image   string
	BuildID string
	WarehouseURL string
	WarehouseToken string
	box     boxfile.Boxfile
}

//
func (sync *Sync) Run() error {

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// set the boxfile and make sure its valid
	if err := sync.setBoxfile(); err != nil {
		return err
	}

	// iterate over the code nodes and build containers for each of them
	for _, codeName := range sync.box.Nodes("code") {
		// pull the image from the boxfile; default to a reasonable alternative if
		// none is given
		image := sync.box.Node(codeName).StringValue("image")
		if image == "" {
			image = "nanobox/code:v1"
		}

		// create a new process config for code ensuring it has access to the warehouse
		// and the boxfile
		codeSetup := Setup{
			App: sync.App,
			Name: codeName,
			Image: image,
		}
		// run the code setup process with the new config
		err := codeSetup.Run()
		if err != nil {
			return fmt.Errorf("code_setup (%s): %s\n", codeName, err.Error())
		}

		codeConfigure := Configure{
			App: sync.App,
			Component: codeSetup.Component,
			BuildID: sync.BuildID,
			WarehouseURL: sync.WarehouseURL,
			WarehouseToken: sync.WarehouseToken,
		}
		// configure this code container
		err = codeConfigure.Run()
		if err != nil {
			return fmt.Errorf("code_configure (%s): %s\n", codeName, err.Error())
		}
	}

	return nil
}

// setBoxfile ...
func (sync *Sync) setBoxfile() error {
	sync.box = boxfile.New([]byte(sync.App.DeployedBoxfile))
	if !sync.box.Valid {
		return fmt.Errorf("Invalid Boxfile")
	}

	return nil
}
