package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processLinkList ...
type processLinkList struct {
	control ProcessControl
}

//
func init() {
	Register("link_list", linkFunc)
}

//
func linkFunc(conf ProcessControl) (Processor, error) {
	return processLinkList{conf}, nil
}

//
func (linkList processLinkList) Results() ProcessControl {
	return linkList.control
}

//
func (linkList processLinkList) Process() error {

	// store the auth token
	links := models.AppLinks{}
	err := data.Get(config.AppName()+"_meta", "links", &links)
	fmt.Printf("%+v\n", links)

	return err
}
