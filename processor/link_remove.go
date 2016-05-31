package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type linkRemove struct {
	control ProcessControl
}

func init() {
	Register("link_remove", linkRemoveFunc)
}

func linkRemoveFunc(conf ProcessControl) (Processor, error) {
	return linkRemove{conf}, nil
}

func (self linkRemove) Results() ProcessControl {
	return self.control
}

func (self linkRemove) Process() error {
	links := models.AppLinks{}
	data.Get(util.AppName()+"_meta", "links", &links)
	delete(links, self.control.Meta["alias"])
	return data.Put(util.AppName()+"_meta", "links", links)
}
