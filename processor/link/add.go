package link

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/odin"
)

// processLinkAdd
type processLinkAdd struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("link_add", linkAddFn)
}

//
func linkAddFn(conf processor.ProcessControl) (processor.Processor, error) {
	return processLinkAdd{conf}, nil
}

//
func (linkAdd processLinkAdd) Results() processor.ProcessControl {
	return linkAdd.control
}

//
func (linkAdd processLinkAdd) Process() error {
	if linkAdd.control.Meta["app"] == "" {

		fmt.Println("you need to provide a app name to link to")
		os.Exit(1)
	}

	// get app id
	app, err := odin.App(linkAdd.control.Meta["app"])
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
