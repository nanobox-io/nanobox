package processor

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/production_api"
)

type linkAdd struct {
	config ProcessConfig
}

func init() {
	Register("link_add", linkAddFunc)
}

func linkAddFunc(conf ProcessConfig) (Processor, error) {
	return linkAdd{conf}, nil
}

func (self linkAdd) Results() ProcessConfig {
	return self.config
}

func (self linkAdd) Process() error {
	if self.config.Meta["name"] == "" {
		fmt.Println("you need to provide a app name to link to")
		os.Exit(1)
	}

	// get app id
	app, err := production_api.App(self.config.Meta["name"])
	if err != nil {
		return err
	}
	if self.config.Meta["alias"] == "" {
		self.config.Meta["alias"] = "default"
	}
	// store the auth token
	link := models.AppLinks{}
	data.Get(util.AppName()+"_meta", "links", &link)
	link[self.config.Meta["alias"]] = app.ID
	return data.Put(util.AppName()+"_meta", "links", link)
}
