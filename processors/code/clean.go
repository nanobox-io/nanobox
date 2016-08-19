package code

import (
	"github.com/jcelliott/lumber"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Clean ...
type Clean struct {
	App models.App
}

// clean all the code services after a dev deploy; unlike the service clean it
// doenst clean ones no longer in the box file but instead removes them all.
func (clean *Clean) Run() error {

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	components, err := models.AllComponentsByApp(clean.App.ID)
	if err != nil {
		lumber.Error("code:Clean:models.AllComponentsByApp(%s): %s", clean.App.ID, err.Error())
		return err
	}

	// get all the code services and remove them
	for _, component := range components {

		// only destroy code type containers
		if component.Type == "code" {

			// run a code destroy
			codeDestroy := Destroy{
				Component: component,
			}

			if err := codeDestroy.Run(); err != nil {
				// we probably dont wnat to break the process just try the rest
				// and report the errors.
			}
		}
	}

	return nil
}
