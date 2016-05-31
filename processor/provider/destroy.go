package provider

import (
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

type providerDestroy struct {
	control processor.ProcessControl
}

func providerDestroyFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerDestroy{control}, nil
}

func (self providerDestroy) Results() processor.ProcessControl {
	return self.control
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
