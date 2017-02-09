package remote

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// Remove ...
func Remove(envModel *models.Env, alias string) error {

	delete(envModel.Remotes, alias)

	// persist the model
	if err := envModel.Save(); err != nil {
		return util.ErrorAppend(err, "failed to remove remote")
	}

	fmt.Printf("\n%s %s remote removed\n\n", display.TaskComplete, alias)

	return nil
}
