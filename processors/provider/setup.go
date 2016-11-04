package provider

import (
	"fmt"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Setup sets up the provider (launch VM, etc)
func Setup() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Starting Nanobox")
	defer display.CloseContext()

	if provider.IsReady() {

		display.StartTask("Skipping (already running)")
		display.StopTask()

		// initialize the docker client
		if err := Init(); err != nil {
			return fmt.Errorf("failed to initialize docker for provider: %s", err.Error())
		}

		return nil
	}

	// install the provider (VM)
	if err := util.Retry(provider.Install, 2, 10*time.Second); err != nil {
		lumber.Error("provider:Setup:provider.Install(): %s", err.Error())
		return fmt.Errorf("failed to install the provider: %s", err.Error())
	}

	// create the provider (VM)
	if err := util.Retry(provider.Create, 2, 10*time.Second); err != nil {
		lumber.Error("provider:Setup:provider.Create(): %s", err.Error())
		return fmt.Errorf("failed to create the provider: %s", err.Error())
	}

	// start the provider (VM)
	if err := util.Retry(provider.Start, 2, 10*time.Second); err != nil {
		lumber.Error("provider:Setup:provider.Start(): %s", err.Error())
		return fmt.Errorf("failed to start the provider: %s", err.Error())
	}

	// fetch the provider model
	providerModel, _ := models.LoadProvider()

	display.StartTask("Joining virtual network")

	// attach the network to the host stack
	if err := setupNetwork(providerModel); err != nil {
		return fmt.Errorf("failed to setup the provider network: %s", err.Error())
	}

	// attach the network to the host stack
	if err := setDefaultIP(providerModel); err != nil {
		return fmt.Errorf("failed to setup the provider network: %s", err.Error())
	}

	display.StopTask()

	// initialize the docker client
	if err := Init(); err != nil {
		return fmt.Errorf("failed to initialize docker for provider: %s", err.Error())
	}

	return nil
}

// setupNetwork sets up the provider network
func setupNetwork(providerModel *models.Provider) error {
	// short-circuit if this is already done
	if providerModel.HostIP != "" {
		return nil
	}

	// reserve an IP to be used for mounting
	mountIP, err := dhcp.ReserveGlobal()
	if err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setupNetwork:dhcp.ReserveGlobal(): %s", err.Error())
		return fmt.Errorf("failed to reserve a global IP: %s", err.Error())
	}
	providerModel.MountIP = mountIP.String()

	// retrieve the provider's Host IP
	hostIP, err := provider.HostIP()
	if err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setupNetwork:provider.HostIP(): %s", err.Error())
		return fmt.Errorf("unable to retrieve the host IP from the provider: %s", err.Error())
	}
	providerModel.HostIP = hostIP

	// persist the IPs for later use
	if err := providerModel.Save(); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to persist the provider model: %s", err.Error())
	}

	return nil
}

// set the default ip everytime
func setDefaultIP(providerModel *models.Provider) error {

	// add the mount IP to the provider
	if err := provider.AddIP(providerModel.MountIP); err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setDefaultIP:provider.AddIP(%s): %s", providerModel.MountIP, err.Error())
		return fmt.Errorf("failed to add an IP to the provider for mounting: %s", err.Error())
	}

	// set the mount IP as the default gateway
	if err := provider.SetDefaultIP(providerModel.MountIP); err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setDefaultIP:provider.SetDefaultIP(%s): %s", providerModel.MountIP, err.Error())
		return fmt.Errorf("failed to set the mount IP as the default gateway: %s", err.Error())
	}

	return nil
}
