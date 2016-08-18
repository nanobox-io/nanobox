package app

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/app/dns"
	"github.com/nanobox-io/nanobox/processor/component"

	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Destroy ...
type Destroy struct {
	App models.App
}

//
func (destroy *Destroy) Run() error {

	if destroy.App.ID == "" {
		// the app doesnt exist
		return nil
	}

	dockerInit()

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	// remove the dev container if there is one
	// but dont catch any errors because there
	// may not be a container
	docker.ContainerRemove(fmt.Sprintf("nanobox_%s", destroy.App.ID))

	// remove all app components
	if err := destroy.removeComponents(); err != nil {
		return err
	}

	// release my ips
	if err := destroy.releaseIPs(); err != nil {
		return err
	}

	// remove all dns entries for this app
	dnsRemoveAll := dns.RemoveAll{destroy.App}
	if err := dnsRemoveAll.Run(); err != nil {
		// report the error but dont stop the process
	}

	// destroy the app
	if err := destroy.App.Delete(); err != nil {
		lumber.Error("app:Destroy:app.Delete(): %s", err.Error())
		return err
	}

	return nil
}

// releaseIPs releases necessary app-global ip addresses
func (destroy *Destroy) releaseIPs() error {

	// release all of the external IPs
	for _, ip := range destroy.App.GlobalIPs {
		// release the IP
		if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
			lumber.Error("app:Destroy:releaseIPs:dhcp.ReturnIP(%s): %s", ip, err.Error())
			return err
		}
	}

	// release all of the local IPs
	for _, ip := range destroy.App.LocalIPs {
		// release the IP
		if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
			lumber.Error("app:Destroy:releaseIPs:dhcp.ReturnIP(%s): %s", ip, err.Error())
			return err
		}
	}

	return nil
}

// removeComponents gets all the components in the app and remove them
func (destroy Destroy) removeComponents() error {

	components, err := models.AllComponentsByApp(destroy.App.ID)
	if err != nil {
		lumber.Error("app:Destroy:removeComponents:models.AllComponentsByApp(%s) %s", destroy.App.ID, err.Error())
		return fmt.Errorf("unable to retrieve components: %s", err.Error())
	}

	for _, comp := range components {

		// do not remove a inprogress build
		if comp.Name == "build" {
			continue
		}

		// creat a component destroy
		componentDestroy := component.Destroy{
			App:       destroy.App,
			Component: comp,
		}

		// run it
		if err := componentDestroy.Run(); err != nil {
			// continue but report the error
		}

	}

	return nil
}
