package processor

import (
	"github.com/nanobox-io/nanobox/processor/provider"
)

// Destroy ...
type Destroy struct {
}

//
func (destroy Destroy) Run() error {
	providerDestroy := provider.Destroy{}
	// run a provider destroy
	return providerDestroy.Run()
}
