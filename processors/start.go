package processors

import (
	"github.com/nanobox-io/nanobox/processors/provider"
)

type Start struct {}

// Run starts the provider (VM)
func (start Start) Run() error {
	// run a provider setup
	providerSetup := provider.Setup{}
	return providerSetup.Run()
}
