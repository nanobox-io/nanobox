package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processProviderStop ...
type processProviderStop struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("provider_stop", providerStopFn)
}

//
func providerStopFn(control processor.ProcessControl) (processor.Processor, error) {
	return processProviderStop{control}, nil
}

//
func (providerStop processProviderStop) Results() processor.ProcessControl {
	return providerStop.control
}

//
func (providerStop processProviderStop) Process() error {
	// establish a global lock to ensure we're the only ones bringing down
	// the provider. Also we need to ensure that we release the lock even
	// if we error out.
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	return provider.Stop()
}
