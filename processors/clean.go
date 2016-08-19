package processors

import (
	"os"

	"github.com/nanobox-io/nanobox/util/locker"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
)

// Clean ...
type Clean struct {
}

//
func (clean *Clean) Run() error {

	// aquire a global lock because we are going to be removing several apps
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	// collect all the apps
	envs, err := models.AllEnvs()
	if err != nil {
		return err
	}

	// check to see if the app folder still exists
	for _, e := range envs {

		if !folderExists(e.Directory) {

			// create an environment destroy for defunct environment
			destroy := env.Destroy{Env: e}

			if err := destroy.Run(); err != nil {
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
