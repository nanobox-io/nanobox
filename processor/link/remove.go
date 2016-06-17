package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processLinkRemove ...
type processLinkRemove struct {
	control processor.ProcessControl
	alias   string
}

//
func init() {
	processor.Register("link_remove", linkRemoveFn)
}

//
func linkRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	linkRemove := &processLinkRemove{control: control}
	return linkRemove, linkRemove.validateMeta()
}

//
func (linkRemove processLinkRemove) Results() processor.ProcessControl {
	return linkRemove.control
}

//
func (linkRemove processLinkRemove) Process() error {

	//
	links := models.AppLinks{}
	if err := data.Get(config.AppName()+"_meta", "links", &links); err != nil {
		//
	}

	//
	delete(links, linkRemove.alias)

	return data.Put(config.AppName()+"_meta", "links", links)
}

// validateMeta validates that the required metadata exists
func (linkRemove *processLinkRemove) validateMeta() error {

	// set alias (required) and ensure it's provided
	linkRemove.alias = linkRemove.control.Meta["alias"]
	if linkRemove.alias == "" {
		return fmt.Errorf("Missing required meta value 'alias'")
	}

	return nil
}
