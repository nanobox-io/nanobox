package provider

import (
	"github.com/nanobox-io/nanobox/processor"
)

type providerSetup struct {
	config processor.ProcessConfig
}


func providerSetupFunc(config processor.ProcessConfig) (Sequence, error) {
	// confirm the provider is an accessable one that we support.

	return providerSetup{config}, nil
}


func (self providerSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self providerSetup) Process() error {
	// TODO: Braxton...
}