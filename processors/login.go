package processors

import (
	"fmt"


	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Process ...
func Login(username, password, endpoint string) error {

	// request Username/Password if missing
	if username == "" {
		user, err := display.ReadUsername()
		if err != nil {
			return fmt.Errorf("unable to retrieve username: %s", err)
		}
		username = user
	}

	if password == "" {
		// ReadPassword prints Password: already
		pass, err := display.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %s", err.Error())
		}
		password = pass
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
