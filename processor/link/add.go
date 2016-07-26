package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/odin"
)

// processLinkAdd
type processLinkAdd struct {
	control processor.ProcessControl
	app     string
	alias   string
}

func init() {
	processor.Register("link_add", linkAddFn)
}

//
func linkAddFn(control processor.ProcessControl) (processor.Processor, error) {
	linkAdd := &processLinkAdd{control: control}
	return linkAdd, linkAdd.validateMeta()
}

func (linkAdd processLinkAdd) Results() processor.ProcessControl {
	return linkAdd.control
}

//
func (linkAdd processLinkAdd) Process() error {

	// get app id
	app, err := odin.App(linkAdd.app)
	if err != nil {
		return err
	}

	// store the auth token
	link := models.AppLinks{}
	if err := data.Get(config.AppID()+"_meta", "links", &link); err != nil {
		//
	}

	//
	link[linkAdd.alias] = app.ID

	return data.Put(config.AppID()+"_meta", "links", link)
}

// validateMeta validates that the required metadata exists
func (linkAdd *processLinkAdd) validateMeta() error {

	// set app (required) and ensure it's provided
	linkAdd.app = linkAdd.control.Meta["app"]
	if linkAdd.app == "" {
		return fmt.Errorf("Missing required meta value 'app'")
	}

	// set alias; if it's not provided set the alias to "default"
	linkAdd.alias = linkAdd.control.Meta["alias"]
	if linkAdd.alias == "" {
		linkAdd.alias = "default"
	}

	return nil
}
