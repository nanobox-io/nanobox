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
func providerStopFunc(control processor.ProcessControl) (processor.Processor, error) {
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
