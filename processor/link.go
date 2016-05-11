package processor

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/production_api"
)

type link struct {
	config ProcessConfig
}

func init() {
	Register("link", linkFunc)
}

func linkFunc(conf ProcessConfig) (Processor, error) {
	return link{conf}, nil
}

func (self link) Results() ProcessConfig {
	return self.config
}

func (self link) Process() error {
	if self.config.Meta["name"] == "" {
		fmt.Println("you need to provide a app name to link to")
		os.Exit(1)
	}

	// get app id
	app, err := production_api.App(self.config.Meta["name"])
	if err != nil {
		return err
	}

	// store the auth token
	link := models.AppLinks{}
	data.Get(util.AppName()"_meta", "links", &link)
	link[self.config.Meta["alias"]] = app.ID
	return data.Put(util.AppName()"_meta", "links", link)
}
