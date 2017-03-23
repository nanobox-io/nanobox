package processors

import (
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/processors/server"
)

// Start starts the provider (VM)
func Start() error {
	// start the nanobox server
	if err := server.Setup(); err != nil {
		return err
	}

	// run a provider setup
	return provider.Setup()
}
