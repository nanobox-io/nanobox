package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/production_api"
)

type login struct {
	control ProcessControl
}

func init() {
	Register("login", loginFunc)
}

func loginFunc(conf ProcessControl) (Processor, error) {
	return login{conf}, nil
}

func (self login) Results() ProcessControl {
	return self.control
}

func (self login) Process() error {
	// request username and password
	if self.control.Meta["username"] == "" {
		self.control.Meta["username"] = print.Prompt("Username:")
	}

	if self.control.Meta["password"] == "" {
		self.control.Meta["password"] = print.Password("Password:")
	}
	// ask odin to verify
	token, err := production_api.Auth(self.control.Meta["username"], self.control.Meta["password"])
	if err != nil {
		return err
	}

	// store the auth token
	auth := models.Auth{token}
	return data.Put("global", "user", auth)
}
