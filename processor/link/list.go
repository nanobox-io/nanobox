package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processLinkList ...
type processLinkList struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("link_list", linkListFn)
}

//
func linkListFn(conf processor.ProcessControl) (processor.Processor, error) {
	return processLinkList{conf}, nil
}

//
func (linkList processLinkList) Results() processor.ProcessControl {
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
