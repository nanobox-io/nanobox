package helpers

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Looks up an app ID on odin and valids it's existence
func OdinAppIDByAlias(alias string) (string, error) {
	// find the app id
	appID := models.AppIDByAlias(alias)

	// if a link doesn't exist, we need to try to look up directly
	if appID == "" {

		appName := alias

		// if an app wasn't specified, we should set the app to the dirname
		if appName == "default" {
			appName = config.AppName()
		}

		// now let's lookup the app on odin
		app, err := ValidateOdinApp(appName)
		if err != nil {
			return "", err
		}

		// set the appID from odin
		appID = app.ID
	}

	// now let's validate the app id in odin
	_, err := ValidateOdinApp(appID)
	if err != nil {
		return "", err
	}

	return appID, nil
}

// Validates an app exists and is accessible on odin
func ValidateOdinApp(slug string) (models.App, error) {
	// fetch the app
	app, err := odin.App(slug)

	// handle errors
	if err != nil {

		if err.Error() == "Unauthorized" {
			fmt.Printf("\n! Sorry, but you don't have access to %s\n\n", slug)
			return app, fmt.Errorf("Unauthorized access to app '%s': %s", slug, err.Error())
		}

		if err.Error() == "Not Found" {
			fmt.Printf("\n! Sorry, the app '%s' doesn't exist\n\n", slug)
			return app, fmt.Errorf("Unknown app '%s': %s", slug, err.Error())
		}

		// All other scenarios
		fmt.Printf("\n! Oops, nanobox is temporarily unreachable. Try again in in just a bit.\n\n")
		return app, fmt.Errorf("Failed to communicate with nanobox: %s", err.Error())
	}

	return app, nil
}
