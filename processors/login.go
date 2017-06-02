package processors

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
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
			return util.ErrorAppend(err, "unable to retrieve username")
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
			return util.ErrorAppend(err, "failed to read password")
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
		fmt.Println(`! The username/password was incorrect, but we're continuing on.
  To reattempt authentication, run 'nanobox login'.
`)
		return nil
	}

	// store the user token
	auth := models.Auth{
		Endpoint: endpoint,
		Key:      token,
	}
	if auth.Save() != nil {
		return util.Errorf("unable to save user authentication")
	}

	display.LoginComplete()

	return nil
}
