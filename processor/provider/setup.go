package provider

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/models"

)

// processProviderSetup ...
type processProviderSetup struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("provider_setup", providerSetupFn)
}

//
func providerSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	return processProviderSetup{control}, nil
}

//
func (providerSetup processProviderSetup) Results() processor.ProcessControl {
	return providerSetup.control
}

//
func (providerSetup processProviderSetup) Process() error {

	// set the provider display level
	provider.Display(!processor.DefaultControl.Quiet)

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

	if err := providerSetup.saveProvider(); err != nil {
		lumber.Error("saveProvider()", err)
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

func (providerSetup processProviderSetup) saveProvider() error {
	mProvider := models.Provider{}
	data.Get("global", "provider", &mProvider)
	
	// if it has already been saved the exit early
	if mProvider.HostIP != "" {
		return nil
	}

	// get a new ip we will use for mounting
	ip, err := dhcp.ReserveGlobal()
	if err != nil {
		return err
	}

	// retrieve the Hosts known ip
	hIP, err := provider.HostIP()
	if err != nil {
		return err
	}
	mProvider.HostIP = hIP
	mProvider.MountIP = ip.String()
	
	return data.Put("global", "provider", mProvider)
}
