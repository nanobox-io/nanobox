package service

import (
	"fmt"
	"net"
	"strings"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type serviceDestroy struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("service_destroy", serviceDestroyFunc)
}

func serviceDestroyFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// if control.Meta["name"] == "" {
	// 	return nil, errors.New("missing image or name")
	// }
	return serviceDestroy{control: control}, nil
}

func (destroy serviceDestroy) Results() processor.ProcessControl {
	return destroy.control
}

func (destroy serviceDestroy) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), destroy.control.Meta["name"], &service)
	if err != nil {
		// cant find service
		return err
	}

	destroy.control.Display(stylish.Bullet("Destroying %s", destroy.control.Meta["name"]))

	err = docker.ContainerRemove(fmt.Sprintf("nanobox-%s-%s", util.AppName(), destroy.control.Meta["name"]))
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

	err = ip_control.ReturnIP(net.ParseIP(service.ExternalIP))
	if err != nil {
		return err
	}

	err = ip_control.ReturnIP(net.ParseIP(service.InternalIP))
	if err != nil {
		return err
	}

	if err := destroy.removeEnvVars(); err != nil {
		return err
	}

	// save the service
	return data.Delete(util.AppName(), destroy.control.Meta["name"])
}

// removeEnvVars removes any env vars associated with this service
func (destroy serviceDestroy) removeEnvVars() error {
	// fetch the environment variables model
	envVars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &envVars)

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	name := destroy.control.Meta["name"]
	prefix := strings.ToUpper(strings.Replace(name, ".", "_", -1))

	// we loop over all environment variables and see if the key contains
	// the prefix above. If so, we delete the item.
	for key, _ := range envVars {
		if strings.HasPrefix(key, prefix) {
			delete(envVars, key)
		}
	}

	// persist the evars
	if err := data.Put(util.AppName()+"_meta", "env", envVars); err != nil {
		return err
	}

	return nil
}
