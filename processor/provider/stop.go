package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
)

type providerStop struct {
	control processor.ProcessControl
}

func providerStopFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerStop{control}, nil
}

func (self providerStop) Results() processor.ProcessControl {
	return self.control
}

func (self providerStop) Process() error {
	return provider.Stop()
}
