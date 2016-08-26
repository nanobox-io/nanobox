package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/locker"
)

//
func Clean(envModels []*models.Env) error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	for _, envModel := range envModels {
		// check to see if the app folder still exists
		if !util.FolderExists(envModel.Directory) {

			if err := env.Destroy(envModel); err != nil {
				return fmt.Errorf("unable to destroy environment(%s): %s", envModel.Name, err)
			}
		}
	}

	return nil
}
