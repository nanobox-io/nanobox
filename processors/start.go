package processor

import (
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
)

// Start ...
type Start struct {
}

//
func (start Start) Run() error {
	display.OpenContext("start provider")
	defer display.CloseContext()

	// run a provider setup
	providerStart := provider.Setup{}
	return providerStart.Run()
}
