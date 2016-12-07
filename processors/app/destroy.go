package app

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app/dns"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Destroy removes the app from the provider and the database
func Destroy(appModel *models.App) error {
	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if this app isn't created
	if appModel.IsNew() {
		return nil
	}

	// load the env for the display context
	envModel, err := appModel.Env()
	if err != nil {
		lumber.Error("app:Start:models.App.Env(): %s", err.Error())
		return fmt.Errorf("failed to load app env: %s", err.Error())
	}

	if err := dns.RemoveAll(appModel); err != nil {
		return fmt.Errorf("failed to remove dns aliases")
	}

	display.OpenContext("%s (%s)", envModel.Name, appModel.DisplayName())
	defer display.CloseContext()

	// remove the dev container if there is one
	docker.ContainerRemove(fmt.Sprintf("nanobox_%s", appModel.ID))

	// destroy the associated components
	if err := destroyComponents(appModel); err != nil {
		return fmt.Errorf("failed to destroy components: %s", err.Error())
	}

	// release IPs
	if err := releaseIPs(appModel); err != nil {
		return fmt.Errorf("failed to release IPs: %s", err.Error())
	}

	// destroy the app model
	if err := appModel.Delete(); err != nil {
		lumber.Error("app:Destroy:models.App{ID:%s}.Destroy(): %s", appModel.ID, err.Error())
		return fmt.Errorf("failed to delete app model: %s", err.Error())
	}

	return nil
}

// destroyComponents destroys all the components of this app
func destroyComponents(appModel *models.App) error {
	display.OpenContext("Removing components")
	defer display.CloseContext()

	componentModels, err := appModel.Components()
	if err != nil {
		lumber.Error("app:destroyComponents:models.App{ID:%s}.Components() %s", appModel.ID, err.Error())
		return fmt.Errorf("unable to retrieve components: %s", err.Error())
	}

	if len(componentModels) == 0 {
		display.StartTask("Skipping (no components)")
		display.StopTask()
		return nil
	}

	for _, componentModel := range componentModels {
		if err := component.Destroy(appModel, componentModel); err != nil {
			return fmt.Errorf("failed to destroy app component: %s", err.Error())
		}
	}

	return nil
}

// releaseIPs releases the app-level ip addresses
func releaseIPs(appModel *models.App) error {
	display.StartTask("Releasing IPs")
	defer display.StopTask()

	// release all of the local IPs
	for _, ip := range appModel.LocalIPs {
		// release the IP
		if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
			display.ErrorTask()
			lumber.Error("app:Destroy:releaseIPs:dhcp.ReturnIP(%s): %s", ip, err.Error())
			return fmt.Errorf("failed to release IP: %s", err.Error())
		}
	}

	return nil
}
