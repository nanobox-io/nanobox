package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

//
func Clean(envModels []*models.Env) error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Cleaning stale environments")
	defer display.CloseContext()

	// if any of the apps are stale, we'll mark this to true
	stale := false

	for _, envModel := range envModels {
		// check to see if the app folder still exists
		if !util.FolderExists(envModel.Directory) {

			if err := env.Destroy(envModel); err != nil {
				return fmt.Errorf("unable to destroy environment(%s): %s", envModel.Name, err)
			}
		}
	}
	
	if !stale {
		display.StartTask("Skipping (none detected)")
		display.StopTask()
	}

	return nil
}
