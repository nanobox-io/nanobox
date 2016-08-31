package processors

import (
	"github.com/nanobox-io/nanobox/processors/provider"
	// "github.com/nanobox-io/nanobox/util/display"
)

// Start starts the provider (VM)
func Start() error {
	// display.OpenContext("Starting Nanobox")
	// defer display.CloseContext()
	
	// run a provider setup
	return provider.Setup()
}
