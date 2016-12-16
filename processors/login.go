package processors

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Process ...
func Login(username, password, endpoint string) error {

	// request Username/Password if missing
	if username == "" && os.Getenv("NANOBOX_USERNAME") != "" {
		username = os.Getenv("NANOBOX_USERNAME")
	}
	if username == "" {
		user, err := display.ReadUsername()
		if err != nil {
			return fmt.Errorf("unable to retrieve username: %s", err)
		}
		username = user
	}

	if password == "" && os.Getenv("NANOBOX_PASSWORD") != "" {
		password = os.Getenv("NANOBOX_PASSWORD")
	}

	if password == "" {
		// ReadPassword prints Password: already
		pass, err := display.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %s", err.Error())
		}
		password = pass
	}

	if endpoint == "" && os.Getenv("NANOBOX_ENDPOINT") != "" {
		endpoint = os.Getenv("NANOBOX_ENDPOINT")
	}

	if endpoint == "" {
		endpoint = "nanobox"
	}

	// set the odin endpoint
	odin.SetEndpoint(endpoint)

	// verify that the user exists
	token, err := odin.Auth(username, password)
	if err != nil {
		fmt.Println("! Incorrect username or password")
		return nil
	}

	// store the user token
	auth := models.Auth{
		Endpoint: endpoint,
		Key:      token,
	}
	if auth.Save() != nil {
		return fmt.Errorf("unable to save user authentication")
	}

	fmt.Printf("%s You've successfully logged in\n", display.TaskComplete)

	return nil
}
