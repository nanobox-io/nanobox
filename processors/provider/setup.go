package provider

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/display"
)

type Setup struct {}

func (setup Setup) Run() error {
	
	display.StartTask("preparing provider")

	// ensure we have an exclusive lock while working with the provider
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	// create the provider (VM)
	if err := provider.Create(); err != nil {
		display.ErrorTask()
		return err
	}

	display.StopTask()
	display.StartTask("booting provider")

	// start the provider (VM)
	if err := provider.Start(); err != nil {
		display.ErrorTask()
		return err
	}

	// attach the network to the host stack
	if err := setup.setupNetwork(); err != nil {
		display.ErrorTask()
		return err
	}

	// fetch the docker env from the provider
	if err := provider.DockerEnv(); err != nil {
		display.ErrorTask()
		return err
	}

	// setup the docker client with the docker environment
	if err := docker.Initialize("env"); err != nil {
		display.ErrorTask()
		return err
	}

	display.StopTask()
	
	return nil
}

// setupNetwork attaches the provider network to the host stack
func (setup Setup) setupNetwork() error {
	// fetch the provider model
	model, _ := models.LoadProvider()
	
	// short-circuit if this is already done
	if model.HostIP != "" {
		return nil
	}
	
	// reserve an IP to be used for mounting
	mountIP, err := dhcp.ReserveGlobal()
	if err != nil {
		lumber.Error("provider:Setup:setupNetwork:dhcp.ReserveGlobal(): %s", err.Error())
		return err
	}
	
	// fetch the host ip from the provider
	hostIP, err := provider.HostIP()
	if err != nil {
		lumber.Error("provider:Setup:setupNetwork:provider.HostIP(): %s", err.Error())
		return err
	}
	
	// persist the IPs for later use
	model.MountIP = mountIP
	model.HostIP  = hostIP
	if err := model.Save(); err != nil {
		return err
	}
	
	// let's attach the mount IP to the provider
	if err := provider.AddIP(mountIP); err != nil {
		lumber.Error("provider:Setup:setupNetwork:provider.AddIP(%s): %s", mountIP, err.Error())
		return err
	}
	
	// now let's set the mount IP as the default route
	if err := provider.SetDefaultIP(mountIP); err != nil {
		lumber.Error("provider:Setup:setupNetwork:provider.SetDefaultIP(%s): %s", mountIP, err.Error())
		return err
	}
	
	return nil
}
