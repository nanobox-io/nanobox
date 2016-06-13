package processor

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/productionAPI"
)

// processLinkAdd
type processLinkAdd struct {
	control ProcessControl
}

//
func init() {
	Register("link_add", linkAddFunc)
}

//
func linkAddFunc(conf ProcessControl) (Processor, error) {
	return processLinkAdd{conf}, nil
}

//
func (linkAdd processLinkAdd) Results() ProcessControl {
	return linkAdd.control
}

//
func (linkAdd processLinkAdd) Process() error {
	if linkAdd.control.Meta["name"] == "" {

		fmt.Println("you need to provide a app name to link to")
		os.Exit(1)
	}

	// get app id
	app, err := productionAPI.App(linkAdd.control.Meta["name"])
	if err != nil {
		return err
	}

	//
	if linkAdd.control.Meta["alias"] == "" {
		linkAdd.control.Meta["alias"] = "default"
	}

	// store the auth token
	link := models.AppLinks{}
	data.Get(config.AppName()+"_meta", "links", &link)
	link[linkAdd.control.Meta["alias"]] = app.ID

	return data.Put(config.AppName()+"_meta", "links", link)
}
