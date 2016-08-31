package code

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// clean all the code services after a dev deploy; unlike the service clean it
// doenst clean ones no longer in the box file but instead removes them all.
func Clean(appModel *models.App) error {

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// get all the components
	componentModels, err := appModel.Components()
	if err != nil {
		lumber.Error("code:Clean:models.App{ID:%s}.Components(): %s", appModel.ID, err.Error())
		return err
	}

	// make sure we only show the context one time
	messaged := false

	// remove components that are of the code type
	for _, componentModel := range componentModels {

		// only destroy code type containers
		if componentModel.Type == "code" {
			if messaged {
				messaged = true
				display.OpenContext("cleaning previous code components")
				defer display.CloseContext()
			}

			// run a code destroy
			if err := Destroy(componentModel); err != nil {
				// TODO: error message
				// we probably dont wnat to break the process just try the rest
				// and report the errors.
			}
		}
	}

	return nil
}
