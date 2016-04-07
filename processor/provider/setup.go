package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
)

type providerSetup struct {
	config processor.ProcessConfig
}


func providerSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerSetup{config}, nil
}

func (self providerSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self providerSetup) Process() error {
	err := provider.Create()
	if err != nil {
		return err
	}
	err = provider.Start()
	if err != nil {
		return err
	}
	// mount my folder
	err = provider.AddMount(util.LocalDir(), "/mount/apps/"+util.AppName()"/")
	return nil
}