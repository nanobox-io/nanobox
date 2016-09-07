package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Process ...
func Logout(endpoint string) error {

	if endpoint == "" {
		endpoint = "nanobox"
	}

	// lookup the auth by the endpoint
	auth, _ := models.LoadAuthByEndpoint(endpoint)
	
	// short-circuit if the auth is already deleted
	if auth.IsNew() {
		fmt.Printf("%s Already logged out\n", display.TaskComplete)
		return nil
	}

	// remove token from database
	if err := auth.Delete(); err != nil {
		return fmt.Errorf("failed to delete user authentication: %s", err.Error())
	}

	fmt.Printf("%s You've logged out\n", display.TaskComplete)

	return nil
}
