package processors

import (
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Start starts the provider (VM)
func Start() error {
	// run a provider setup
	return provider.Setup()
}
