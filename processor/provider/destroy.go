package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

type providerDestroy struct {
	config processor.ProcessConfig
}

func providerDestroyFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerDestroy{config}, nil
}

func (self providerDestroy) Results() processor.ProcessConfig {
	return self.config
}

func (self providerDestroy) Process() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	err := self.RemoveDatabase()
	if err != nil {
		return err
	}

	return provider.Destroy()
}
