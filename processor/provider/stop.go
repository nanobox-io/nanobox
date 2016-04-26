package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
)

type providerStop struct {
	config processor.ProcessConfig
}

func providerStopFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerStop{config}, nil
}

func (self providerStop) Results() processor.ProcessConfig {
	return self.config
}

func (self providerStop) Process() error {
	return provider.Stop()
}
