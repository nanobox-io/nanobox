package helpers

import (
	"fmt"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Validates an app exists and is accessible on odin
func ValidateOdinApp(slug string) error {
	// fetch the app
	_, err := odin.App(slug)

	// handle errors
	if err != nil {

		lumber.Error("helpers: ValidateOdinApp(%s): %s", slug, err)

		if strings.Contains(err.Error(), "Unauthorized") {
			fmt.Printf("\n! Sorry, but you don't have access to %s\n\n", slug)
			return util.ErrorAppend(err, "Unauthorized access to app '%s'", slug)
		}

		if strings.Contains(err.Error(), "Not Found") {
			fmt.Printf("\n! Sorry, the app '%s' doesn't exist\n\n", slug)
			return util.ErrorAppend(err, "Unknown app '%s'", slug)
		}

		// All other scenarios
		fmt.Printf("\n%s\n\n", err.Error())
		return util.ErrorAppend(err, "Failed to communicate with nanobox")
	}

	return nil
}
