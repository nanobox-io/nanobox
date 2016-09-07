package processors

import (
	"fmt"

	printutil "github.com/sdomino/go-util/print"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Process ...
func Login(username, password, endpoint string) error {

	// request Username/Password if missing
	if username == "" {
		// add in tylers display system for prompting
		username = printutil.Prompt("Username: ")
	}

	if password == "" {
		// ReadPassword prints Password: already
		pass, err := util.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %s", err.Error())
		}
		password = pass
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
		Key: token,
	}
	if auth.Save() != nil {
		return fmt.Errorf("unable to save user authentication")
	}

	fmt.Printf("%s You've successfully logged in\n", display.TaskComplete)

	return nil
}
