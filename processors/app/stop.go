package app

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	process_provider "github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Stop will stop all services associated with an app
func Stop(appModel *models.App) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is already down
	// TODO: also check if any containers are running
	if appModel.Status != "up" {
		return nil
	}

	// load the env for the display context
	envModel, err := appModel.Env()
	if err != nil {
		lumber.Error("app:Stop:models.App.Env(): %s", err.Error())
		return fmt.Errorf("failed to load app env: %s", err.Error())
	}

	display.OpenContext("%s (%s)", envModel.Name, appModel.Name)
	defer display.CloseContext()

	// initialize docker for the provider
	if err := process_provider.Init(); err != nil {
		return fmt.Errorf("failed to initialize docker environment: %s", err.Error())
	}

	// stop all app components
	if err := component.StopAll(appModel); err != nil {
		return fmt.Errorf("failed to stop all app components: %s", err.Error())
	}

	// stop any dev containers
	stopDevContainer(appModel)

	// set the status to down
	appModel.Status = "down"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Stop:models.App.Save(): %s", err.Error())
		return fmt.Errorf("failed to persist app status: %s", err.Error())
	}

	return nil
}

func stopDevContainer(appModel *models.App) error {
	// grab the container info
	container, err := docker.GetContainer(fmt.Sprintf("nanobox_%s", appModel.ID))
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}

	// remove the container
	if err := docker.ContainerRemove(container.ID); err != nil {
		lumber.Error("dev:console:teardown:docker.ContainerRemove(%s): %s", container.ID, err)
		return fmt.Errorf("failed to remove dev container: %s", err.Error())
	}

	// extract the container IP
	ip := docker.GetIP(container)

	// detach dev container from the network
	if err := detachNetwork(appModel, ip); err != nil {
		return fmt.Errorf("failed to detach dev container from network: %s", err.Error())
	}

	// return the container IP back to the IP pool
	if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
		lumber.Error("dev:console:teardown:dhcp.ReturnIP(%s): %s", ip, err)

		lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		return fmt.Errorf("failed to return unused IP back to pool: %s", err.Error())
	}
	return nil
}

// detachNetwork detaches the container from the host network
func detachNetwork(appModel *models.App, containerIP string) error {

	// fetch the envIP
	envIP := appModel.GlobalIPs["env"]

	// remove nat
	if err := provider.RemoveNat(envIP, containerIP); err != nil {
		lumber.Error("dev:detachNetwork:provider.RemoveNat(%s, %s):%s", envIP, containerIP, err.Error())
		return fmt.Errorf("failed to remove NAT from container: %s", err.Error())
	}

	// remove the IP from the provider
	if err := provider.RemoveIP(envIP); err != nil {
		lumber.Error("dev:detachNetwork:provider.RemoveIP(%s):%s", envIP, err.Error())
		return fmt.Errorf("failed to remove the IP from the provider: %s", err.Error())
	}

	return nil
}
