package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processProviderDestroy ...
type processProviderDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("provider_destroy", providerDestroyFn)
}

//
func providerDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return processProviderDestroy{control}, nil
}

//
func (providerDestroy processProviderDestroy) Results() processor.ProcessControl {
	return providerDestroy.control
}

//
func (providerDestroy processProviderDestroy) Process() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	//
	if err := providerDestroy.removeDatabase(); err != nil {
		return err
	}

	return provider.Destroy()
}
