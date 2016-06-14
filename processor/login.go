package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/odin"

	printutil "github.com/sdomino/go-util/print"
)

type processLogin struct {
	control ProcessControl
}

func init() {
	Register("login", loginFunc)
}

func loginFunc(conf ProcessControl) (Processor, error) {
	return processLogin{conf}, nil
}

// Results ...
func (login processLogin) Results() ProcessControl {
	return login.control
}

// Process ...
func (login processLogin) Process() error {
	// request username and password
	if login.control.Meta["username"] == "" {
		login.control.Meta["username"] = printutil.Prompt("Username:")
	}

	if login.control.Meta["password"] == "" {
		login.control.Meta["password"] = printutil.Password("Password:")
	}
	// ask odin to verify
	token, err := odin.Auth(login.control.Meta["username"], login.control.Meta["password"])
	if err != nil {
		return err
	}

	// store the auth token
	auth := models.Auth{Key: token}
	return data.Put("global", "user", auth)
}
