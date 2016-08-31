package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// Process ...
func Logout() error {

	// remove token from database
	if err := models.DeleteAuth(); err != nil {
		return fmt.Errorf("failed to delete auth: %s", err.Error())
	}

	fmt.Println("Successfully logged out!")

	return nil
}
