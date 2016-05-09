package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/production_api"
)

type login struct {
	config ProcessConfig
}

func init() {
	Register("login", loginFunc)
}

func loginFunc(conf ProcessConfig) (Processor, error) {
	return login{conf}, nil
}

func (self login) Results() ProcessConfig {
	return self.config
}

func (self login) Process() error {
	// request username and password
	if self.config.Meta["username"] == "" {
		self.config.Meta["username"] = print.Prompt("Username:")
	}

	if self.config.Meta["password"] == "" {
		self.config.Meta["password"] = print.Password("Password:")
	}
	// ask odin to verify
	token, err := production_api.Auth(self.config.Meta["username"], self.config.Meta["password"])
	if err != nil {
		return err
	}

	// store the auth token
	auth := models.Auth{token}
	return data.Put("global", "user", auth)
}
