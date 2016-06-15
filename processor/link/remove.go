package link

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processLinkRemove ...
type processLinkRemove struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("link_remove", linkRemoveFn)
}

//
func linkRemoveFn(conf processor.ProcessControl) (processor.Processor, error) {
	return processLinkRemove{conf}, nil
}

//
func (linkRemove processLinkRemove) Results() processor.ProcessControl {
	return linkRemove.control
}

//
func (linkRemove processLinkRemove) Process() error {
	links := models.AppLinks{}
	data.Get(config.AppName()+"_meta", "links", &links)
	delete(links, linkRemove.control.Meta["alias"])
	return data.Put(config.AppName()+"_meta", "links", links)
}
