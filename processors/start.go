package processors

import (
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
)

type Start struct {}

func (start Start) Run() error {
	display.OpenContext("start provider")
	defer display.CloseContext()	
	
	// run the provider setup processor
	return provider.Setup{}.Run()
}
