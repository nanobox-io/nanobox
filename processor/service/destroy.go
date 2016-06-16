package service

import (
	"fmt"
	"net"
	"strings"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

// processServiceDestroy ...
type processServiceDestroy struct {
	control 	processor.ProcessControl
	app				models.App
	service		models.Service
}

//
func init() {
	processor.Register("service_destroy", serviceDestroyFn)
}

//
func serviceDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// if control.Meta["name"] == "" {
	// 	return nil, errors.New("missing image or name")
	// }
	return processServiceDestroy{control: control}, nil
}

//
func (serviceDestroy processServiceDestroy) Results() processor.ProcessControl {
	return serviceDestroy.control
}

//
func (serviceDestroy processServiceDestroy) Process() error {

	if err := serviceDestroy.loadApp(); err != nil {
		return err
	}

	if err := serviceDestroy.loadService(); err != nil {
		return err
	}

	if err := serviceDestroy.printDisplay(); err != nil {
		return err
	}

	if err := serviceDestroy.removeContainer(); err != nil {
		return err
	}

	if err := serviceDestroy.detachNetwork(); err != nil {
		return err
	}

	if err := serviceDestroy.removeEnvVars(); err != nil {
		return err
	}

	if err := serviceDestroy.deleteService(); err != nil {
		return err
	}

	return nil
}

// loadApp loads the app from the database
func (serviceDestroy *processServiceDestroy) loadApp() error {

	// load the app from the database
	if err := data.Get("apps", config.AppName(), &serviceDestroy.app); err != nil {
		return err
	}

	return nil
}

// loadService fetches the service from the database
func (serviceDestroy *processServiceDestroy) loadService() error {
	name := serviceDestroy.control.Meta["name"]

	// the service really shouldn't exist yet, so let's not return the error if it fails
	data.Get(config.AppName(), name, &serviceDestroy.service)

	return nil
}

// printDisplay prints the user display for progress
func (serviceDestroy *processServiceDestroy) printDisplay() error {

	name    := serviceDestroy.control.Meta["name"]
	message := stylish.Bullet("Destroying %s", name)

	// print!
	serviceDestroy.control.Display(message)

	return nil
}

// removeContainer destroys the docker container
func (serviceDestroy *processServiceDestroy) removeContainer() error {

	name := serviceDestroy.control.Meta["name"]
	container := fmt.Sprintf("nanobox-%s-%s", config.AppName(), name)

	if err := docker.ContainerRemove(container); err != nil {
		return err
	}

	return nil
}

// detachNetwork detaches the virtual network from the host
func (serviceDestroy *processServiceDestroy) detachNetwork() error {

	name    := serviceDestroy.control.Meta["name"]
	service := serviceDestroy.service

	if err := provider.RemoveNat(service.ExternalIP, service.InternalIP); err != nil {
		return err
	}

	if err := provider.RemoveIP(service.ExternalIP); err != nil {
		return err
	}

	// don't return the external IP if this is portal
	if name != "portal" {
		if err := dhcp.ReturnIP(net.ParseIP(service.ExternalIP)); err != nil {
			return err
		}
	}


	// don't return the internal IP if it's an app-level cache
	if serviceDestroy.app.LocalIPs[name] != "" {
		if err := dhcp.ReturnIP(net.ParseIP(service.InternalIP)); err != nil {
			return err
		}
	}

	return nil
}

// removeEnvVars removes any env vars associated with this service
func (serviceDestroy processServiceDestroy) removeEnvVars() error {
	// fetch the environment variables model
	envVars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &envVars)

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	name := serviceDestroy.control.Meta["name"]
	prefix := strings.ToUpper(strings.Replace(name, ".", "_", -1))

	// we loop over all environment variables and see if the key contains
	// the prefix above. If so, we delete the item.
	for key := range envVars {
		if strings.HasPrefix(key, prefix) {
			delete(envVars, key)
		}
	}

	// persist the evars
	if err := data.Put(config.AppName()+"_meta", "env", envVars); err != nil {
		return err
	}

	return nil
}

// deleteService deletes the service record from the db
func (serviceDestroy processServiceDestroy) deleteService() error {

	name := serviceDestroy.control.Meta["name"]

	if err := data.Delete(config.AppName(), name); err != nil {
		return err
	}

	return nil
}
