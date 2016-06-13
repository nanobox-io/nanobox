package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processLinkRemove ...
type processLinkRemove struct {
	control ProcessControl
}

//
func init() {
	Register("link_remove", linkRemoveFunc)
}

//
func linkRemoveFunc(conf ProcessControl) (Processor, error) {
	return processLinkRemove{conf}, nil
}

//
func (linkRemove processLinkRemove) Results() ProcessControl {
	return linkRemove.control
}

//
func (linkRemove processLinkRemove) Process() error {
	links := models.AppLinks{}
	data.Get(config.AppName()+"_meta", "links", &links)
	delete(links, linkRemove.control.Meta["alias"])
	return data.Put(config.AppName()+"_meta", "links", links)
}
