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
	control ProcessControl
}

func init() {
	Register("link_add", linkAddFunc)
}

func linkAddFunc(conf ProcessControl) (Processor, error) {
	return linkAdd{conf}, nil
}

func (self linkAdd) Results() ProcessControl {
	return self.control
}

func (self linkAdd) Process() error {
	if self.control.Meta["name"] == "" {
		fmt.Println("you need to provide a app name to link to")
		os.Exit(1)
	}

	// get app id
	app, err := production_api.App(self.control.Meta["name"])
	if err != nil {
		return err
	}
	if self.control.Meta["alias"] == "" {
		self.control.Meta["alias"] = "default"
	}
	// store the auth token
	link := models.AppLinks{}
	data.Get(util.AppName()+"_meta", "links", &link)
	link[self.control.Meta["alias"]] = app.ID
	return data.Put(util.AppName()+"_meta", "links", link)
}
