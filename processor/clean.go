package processor

import (
	"os"
	"strings"

	"github.com/nanobox-io/nanobox/util/locker"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"

)

// processClean ...
type processClean struct {
	control ProcessControl
}

//
func init() {
	Register("clean", cleanFn)
}

//
func cleanFn(control ProcessControl) (Processor, error) {
	return processClean{control}, nil
}

//
func (clean processClean) Results() ProcessControl {
	return clean.control
}

//
func (clean processClean) Process() error {

	// aquire a global lock because we are going to be removing several apps
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	// collect all the apps
	keys, err := data.Keys("apps")
	if err != nil {
		return err
	}

	// check to see if the app folder still exists
	for _, appID := range keys {

		app := models.App{}
		data.Get("apps", appID, &app)

		if !folderExists(app.Directory) {

			// remove apps that no longer exist in the folder
			clean.control.Meta["app_name"] = app.ID
			clean.control.Env = "dev"


			// get the env from the id
			if strings.Contains(app.ID, "_sim") {
				clean.control.Env = "sim"
			}

			err := Run("env_destroy", clean.control)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func folderExists(dirName string) bool {
	dir, err := os.Stat(dirName)
	if err != nil {
		return false
	}
	return dir.IsDir()
}
