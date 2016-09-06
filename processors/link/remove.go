package link

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Remove ...
func Remove(envModel *models.Env, alias string) error {

	delete(envModel.Links, alias)

	// persist the model
	if err := envModel.Save(); err != nil {
		return fmt.Errorf("failed to remove link: %s", err.Error())
	}

	fmt.Printf("\n%s %s link removed\n\n", display.TaskComplete, alias)

	return nil
}
