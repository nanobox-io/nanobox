package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Process ...
func Logout() error {

	// remove token from database
	if err := models.DeleteAuth(); err != nil {
		return fmt.Errorf("failed to delete user authentication: %s", err.Error())
	}

	fmt.Printf("%s You've logged out\n", display.TaskComplete)

	return nil
}
