package provider

import (
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Stop ...
type Stop struct {
}

//
func (providerStop Stop) Run() error {
	// establish a global lock to ensure we're the only ones bringing down
	// the provider. Also we need to ensure that we release the lock even
	// if we error out.
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	return provider.Stop()
}
