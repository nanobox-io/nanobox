package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type link struct {
	config ProcessConfig
}

func init() {
	Register("link_list", linkFunc)
}

func linkFunc(conf ProcessConfig) (Processor, error) {
	return link{conf}, nil
}

func (self link) Results() ProcessConfig {
	return self.config
}

func (self link) Process() error {
	// store the auth token
	links := models.AppLinks{}
	err := data.Get(util.AppName()+"_meta", "links", &links)
	fmt.Printf("%+v\n", links)
	return err
}
