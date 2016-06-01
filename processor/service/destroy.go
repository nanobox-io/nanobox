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

func (self serviceDestroy) Results() processor.ProcessControl {
	return self.control
}

func (self serviceDestroy) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.control.Meta["name"], &service)
	if err != nil {
		// cant find service
		return err
	}

	self.control.Display(stylish.Bullet("Destroying %s", self.control.Meta["name"]))

	err = docker.ContainerRemove(fmt.Sprintf("nanobox-%s-%s", util.AppName(), self.control.Meta["name"]))
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

	if err := self.updateEnvVars(); err != nil {
		return err
	}

	// save the service
	return data.Delete(util.AppName(), self.control.Meta["name"])
}

func (self serviceDestroy) updateEnvVars() error {
	envVars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &envVars)

	envName := strings.ToUpper(strings.Replace(self.control.Meta["name"], ".", "_", -1))
	delete(envVars, envName+"_HOST")
	users := strings.Split(envVars[envName+"_USERS"], " ")
	for _, user := range users {
		delete(envVars, fmt.Sprintf("%s_%s_PW", envName, strings.ToUpper(user)))
	}
	delete(envVars, fmt.Sprintf("%s_USER", envName))
	delete(envVars, fmt.Sprintf("%s_PASS", envName))

	delete(envVars, envName+"_USERS")
	return data.Put(util.AppName()+"_meta", "env", envVars)
}
