package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"

	printutil "github.com/sdomino/go-util/print"
)

type Login struct {
	Username string
	Password string
	token    string
}

// Process ...
func (login Login) Run() error {

	// validate we have all meta information needed
	if err := login.validateUser(); err != nil {
		return err
	}

	// verify that the user exists
	if err := login.verifyUser(); err != nil {
		return err
	}

	fmt.Println("verified user")

	// store the user token
	if err := login.saveUser(); err != nil {
		return err
	}

	return nil
}

// validateUser validates that the required metadata exists
func (login *Login) validateUser() error {

	// request Username/Password if missing
	if login.Username == "" {
		// add in tylers display system for prompting
		login.Username = printutil.Prompt("Username: ")
	}

	if login.Password == "" {
		// ReadPassword prints Password: already
		pass, err := util.ReadPassword()
		if err != nil {
			// TODO: print out the error to the log
		}
		login.Password = pass
	}

	return nil
}

// verifyUser ...
func (login *Login) verifyUser() (err error) {

	//
	if login.token, err = odin.Auth(login.Username, login.Password); err != nil {
		return err
	}

	return nil
}

// saveUser ...
func (login *Login) saveUser() error {

	// store the auth token
	auth := models.Auth{Key: login.token}
	return auth.Save()
}
