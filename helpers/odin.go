package helpers

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/odin"
)

// Validates an app exists and is accessible on odin
func ValidateOdinApp(slug string) error {
	// fetch the app
	_, err := odin.App(slug)
	
	// handle errors
	if err != nil {
		
		lumber.Error("helpers: ValidateOdinApp(%s): %s", slug, err)
		
		if err.Error() == "Unauthorized" {
			fmt.Printf("\n! Sorry, but you don't have access to %s\n\n", slug)
			return fmt.Errorf("Unauthorized access to app '%s': %s", slug, err.Error())
		}

		if err.Error() == "Not Found" {
			fmt.Printf("\n! Sorry, the app '%s' doesn't exist\n\n", slug)
			return fmt.Errorf("Unknown app '%s': %s", slug, err.Error())
		}

		// All other scenarios
		fmt.Printf("\n! Oops, nanobox is temporarily unreachable. Try again in in just a bit.\n\n")
		return fmt.Errorf("Failed to communicate with nanobox: %s", err.Error())
	}

	return nil
}
