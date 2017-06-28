package provider

import (
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider/bridge"
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

	if provider.IsReady() {

		// if we are already ready we may still need to bridge
		display.StartTask("Skipping (already running)")
		display.StopTask()

		if err := Init(); err != nil {
			return util.ErrorAppend(err, "failed to initialize docker for provider")
		}

		if provider.BridgeRequired() {
			if err := bridge.Setup(); err != nil {
				return util.ErrorAppend(err, "failed to setup the network bridge")
			}
		}
		return nil
	}

	display.OpenContext("Starting Nanobox")

	// create the provider (VM)
	if err := util.Retry(provider.Create, 3, 20*time.Second); err != nil {
		lumber.Error("provider:Setup:provider.Create()")
		return util.ErrorAppend(err, "failed to create the provider")
	}

	// start the provider (VM)
	if err := util.Retry(provider.Start, 3, 20*time.Second); err != nil {
		lumber.Error("provider:Setup:provider.Start()")
		return util.ErrorAppend(err, "failed to start the provider")
	}

	// fetch the provider model
	providerModel, _ := models.LoadProvider()
	providerModel.Name = provider.Name()

	display.StartTask("Joining virtual network")

	// attach the network to the host stack
	if err := setupNetwork(providerModel); err != nil {
		return util.ErrorAppend(err, "failed to setup the provider network")
	}

	display.StopTask()

	if err := Init(); err != nil {
		return util.ErrorAppend(err, "failed to initialize docker for provider")
	}

	display.CloseContext()

	if provider.BridgeRequired() {
		if err := bridge.Setup(); err != nil {
			return util.ErrorAppend(err, "failed to setup the network bridge")
		}
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
		lumber.Error("provider:Setup:setupNetwork:dhcp.ReserveGlobal()")
		return util.ErrorAppend(err, "failed to reserve a global IP")
	}
	providerModel.MountIP = mountIP.String()

	// retrieve the provider's Host IP
	hostIP, err := provider.HostIP()
	if err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setupNetwork:provider.HostIP()")
		return util.ErrorAppend(err, "unable to retrieve the host IP from the provider")
	}
	providerModel.HostIP = hostIP

	// persist the IPs for later use
	if err := providerModel.Save(); err != nil {
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to persist the provider model")
	}

	return nil
}

// set the default ip everytime
func setDefaultIP(providerModel *models.Provider) error {

	// add the mount IP to the provider
	if err := provider.AddIP(providerModel.MountIP); err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setDefaultIP:provider.AddIP(%s): %s", providerModel.MountIP, err.Error())
		return util.ErrorAppend(err, "failed to add an IP to the provider for mounting")
	}

	// set the mount IP as the default gateway
	if err := provider.SetDefaultIP(providerModel.MountIP); err != nil {
		display.ErrorTask()
		lumber.Error("provider:Setup:setDefaultIP:provider.SetDefaultIP(%s): %s", providerModel.MountIP, err.Error())
		return util.ErrorAppend(err, "failed to set the mount IP as the default gateway")
	}

	return nil
}
