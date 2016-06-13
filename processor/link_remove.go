package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
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
	data.Get(util.AppName()+"_meta", "links", &links)
	delete(links, linkRemove.control.Meta["alias"])
	return data.Put(util.AppName()+"_meta", "links", links)
}
