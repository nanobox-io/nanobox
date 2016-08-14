package component

import (
	"fmt"
	"net"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

// Destroy ...
type Destroy struct {
	App       models.App
	Component models.Component
}

//
func (destroy *Destroy) Run() error {

	if err := destroy.removeContainer(); err != nil {
		// if im unable to remove the container (especially if it doesnt exist)
		// we shouldnt fail
		lumber.Error("unable to removeContainer: %s", err.Error())
	}

	if err := destroy.removeEvars(); err != nil {
		return fmt.Errorf("unable to removeEvars: %s", err.Error())
	}

	if err := destroy.Component.Delete(); err != nil {
		return fmt.Errorf("unable to deleteComponent: %s", err.Error())
	}

	return nil
}

// removeContainer destroys the docker container
func (destroy *Destroy) removeContainer() error {

	containerName := fmt.Sprintf("nanobox_%s_%s", destroy.App.ID, destroy.Component.Name)

	if err := docker.ContainerRemove(containerName); err != nil {
		return err
	}

	return nil
}

// detachNetwork detaches the virtual network from the host
func (destroy *Destroy) detachNetwork() error {
	name := destroy.Component.Name

	if err := provider.RemoveNat(destroy.Component.ExternalIP, destroy.Component.InternalIP); err != nil {
		return err
	}

	if err := provider.RemoveIP(destroy.Component.ExternalIP); err != nil {
		return err
	}

	// don't return the external IP if this is portal
	if name != "portal" && destroy.App.GlobalIPs[name] == "" {
		if err := dhcp.ReturnIP(net.ParseIP(destroy.Component.ExternalIP)); err != nil {
			return err
		}
	}

	// don't return the internal IP if it's an app-level cache
	if destroy.App.LocalIPs[name] == "" {
		if err := dhcp.ReturnIP(net.ParseIP(destroy.Component.InternalIP)); err != nil {
			return err
		}
	}

	return nil
}

// removeEvars removes any env vars associated with this service
func (destroy Destroy) removeEvars() error {
	// fetch the environment variables
	envVars := destroy.App.Evars

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	prefix := strings.ToUpper(strings.Replace(destroy.Component.Name, ".", "_", -1))

	// we loop over all environment variables and see if the key contains
	// the prefix above. If so, we delete the item.
	for key := range envVars {
		if strings.HasPrefix(key, prefix) {
			delete(envVars, key)
		}
	}

	// persist the evars
	destroy.App.Save()

	return nil
}
