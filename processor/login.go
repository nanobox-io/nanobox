package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/odin"

	printutil "github.com/sdomino/go-util/print"
)

type processLogin struct {
	control  ProcessControl
	username string
	password string
	token    string
}

func init() {
	Register("login", loginFn)
}

func loginFn(control ProcessControl) (Processor, error) {
	return processLogin{control: control}, nil
}

func (login processLogin) Results() ProcessControl {
	return login.control
}

// Process ...
func (login processLogin) Process() error {

	// validate we have all meta information needed
	if err := login.validateMeta(); err != nil {
		return err
	}

	// verify that the user exists
	if err := login.verifyUser(); err != nil {
		return err
	}

	// store the user token
	if err := login.saveUser(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (login *processLogin) validateMeta() error {

	// set username and password
	login.username = login.control.Meta["username"]
	login.password = login.control.Meta["password"]

	// request username/password if missing
	if login.username == "" {
		login.username = printutil.Prompt("Username:")
		fmt.Println("username: ", login.username)
	}

	if login.password == "" {
		// ReadPassword prints Password: already
		if pass, err := util.ReadPassword(); err != nil {
			login.password = pass
		} else {
			login.password = ""
		}
		fmt.Println("password: ", login.password)
	}

	return nil
}

// verifyUser ...
func (login *processLogin) verifyUser() (err error) {

	//
	if login.token, err = odin.Auth(login.username, login.password); err != nil {
		return err
	}

	return nil
}

// saveUser ...
func (login *processLogin) saveUser() error {

	// store the auth token
	if err := data.Put("global", "user", models.Auth{Key: login.token}); err != nil {
		return err
	}

	return nil
}
