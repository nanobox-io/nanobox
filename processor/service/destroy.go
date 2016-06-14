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
	control processor.ProcessControl
}

//
func init() {
	processor.Register("service_destroy", serviceDestroyFunc)
}

//
func serviceDestroyFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// if control.Meta["name"] == "" {
	// 	return nil, errors.New("missing image or name")
	// }
	return processServiceDestroy{control: control}, nil
}

//
func (servicedestroy processServiceDestroy) Results() processor.ProcessControl {
	return servicedestroy.control
}

//
func (servicedestroy processServiceDestroy) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(config.AppName(), servicedestroy.control.Meta["name"], &service)
	if err != nil {
		// cant find service
		return err
	}

	servicedestroy.control.Display(stylish.Bullet("Destroying %s", servicedestroy.control.Meta["name"]))

	err = docker.ContainerRemove(fmt.Sprintf("nanobox-%s-%s", config.AppName(), servicedestroy.control.Meta["name"]))
	if err != nil {
		return err
	}

	err = provider.RemoveNat(service.ExternalIP, service.InternalIP)
	if err != nil {
		return err
	}

	err = provider.RemoveIP(service.ExternalIP)
	if err != nil {
		return err
	}

	err = dhcp.ReturnIP(net.ParseIP(service.ExternalIP))
	if err != nil {
		return err
	}

	err = dhcp.ReturnIP(net.ParseIP(service.InternalIP))
	if err != nil {
		return err
	}

	if err := servicedestroy.removeEnvVars(); err != nil {
		return err
	}

	// save the service
	return data.Delete(config.AppName(), servicedestroy.control.Meta["name"])
}

// removeEnvVars removes any env vars associated with this service
func (servicedestroy processServiceDestroy) removeEnvVars() error {
	// fetch the environment variables model
	envVars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &envVars)

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	name := servicedestroy.control.Meta["name"]
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
