package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
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
	// confirm the provider is an accessable one that we support.
	return processProviderStop{control}, nil
}

//
func (providerStop processProviderStop) Results() processor.ProcessControl {
	return providerStop.control
}

//
func (providerStop processProviderStop) Process() error {
	return provider.Stop()
}
