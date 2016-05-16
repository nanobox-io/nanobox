package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type linkRemove struct {
	config ProcessConfig
}

func init() {
	Register("link_remove", linkRemoveFunc)
}

func linkRemoveFunc(conf ProcessConfig) (Processor, error) {
	return linkRemove{conf}, nil
}

func (self linkRemove) Results() ProcessConfig {
	return self.config
}

func (self linkRemove) Process() error {
	links := models.AppLinks{}
	data.Get(util.AppName()+"_meta", "links", &links)
	delete(links, self.config.Meta["alias"])
	return data.Put(util.AppName()+"_meta", "links", links)
}
