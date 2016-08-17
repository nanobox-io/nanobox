package processor

import (
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/processor/provider"
)

// Destroy ...
type Destroy struct {
}

//
func (destroy Destroy) Run() error {
	display.OpenContext("Destroying nanobox system")
	defer display.CloseContext()
	
	providerDestroy := provider.Destroy{}
	// run a provider destroy
	return providerDestroy.Run()
}
