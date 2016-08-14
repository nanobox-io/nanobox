package processor

import (
	"github.com/nanobox-io/nanobox/processor/provider"
)

// Start ...
type Start struct {
}

//
func (start Start) Run() error {
	// run a provider setup
	providerStart := provider.Setup{}
	return providerStart.Run()
}
