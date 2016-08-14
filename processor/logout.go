package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

type Logout struct {
}

// Process ...
func (logout Logout) Run() error {

	// remove token from database
	if err := models.DeleteAuth(); err != nil {
		return err
	}

	fmt.Println("Successfully logged out!")

	return nil
}
