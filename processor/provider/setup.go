package provider

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processProviderSetup ...
type processProviderSetup struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("provider_setup", providerSetupFunc)
}

//
func providerSetupFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return processProviderSetup{control}, nil
}

//
func (providerSetup processProviderSetup) Results() processor.ProcessControl {
	return providerSetup.control
}

//
func (providerSetup processProviderSetup) Process() error {
	// set the provider display level
	provider.Display(!processor.DefaultConfig.Quiet)

	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if err := provider.Create(); err != nil {
		lumber.Error("Create()", err)
		return err
	}

	if err := provider.Start(); err != nil {
		lumber.Error("Start()", err)
		return err
	}

	if err := provider.DockerEnv(); err != nil {
		lumber.Error("DockerEnv()", err)
		return err
	}

	if err := docker.Initialize("env"); err != nil {
		lumber.Error("docker.Initialize", err)
		return err
	}

	return nil
}
