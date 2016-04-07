package service

import (
	"errors"
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type serviceSetup struct {
	config processor.ProcessConfig
	fail bool
}

var missingImageOrName = errors.New("missing image or name")

func init() {
	processor.Register("serivce_setup", serviceSetupFunc)
	docker.Initialize("env")
}

func serviceSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return serviceSetup{config: config}, nil
}

func (self serviceSetup) clean(fn func()) func() {
	return func() {
		if self.fail {
			fn()
		}
	}
}

func (self serviceSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceSetup) Process() error {
	// make sure i was given a name and image
	if self.config.Meta["name"] == "" || self.config.Meta["image"] == "" {
		return missingImageOrName
	}

	// get the service from the database
	service := models.Service{}
	data.Get(util.AppName(), self.config.Meta["name"], &service)


	// create docker container
	if service.ID != "" {
		// quit early if the service was found to be created already
		return nil
	}

	local_ip, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}
	defer self.clean(func() {
		ip_control.ReturnIP(local_ip)
	})()

	global_ip, err := ip_control.ReserveGlobal()
	if err != nil {
		self.fail = true
		return err
	}
	defer self.clean(func() {
		ip_control.ReturnIP(global_ip)
	})()

	config := docker.ContainerConfig{
		Name: fmt.Sprintf("%s-%s", util.AppName(), self.config.Meta["name"]),
		Image: self.config.Meta["image"],
 		Network: "virt",
 		IP: local_ip.String(),
	}

	container, err := docker.CreateContainer(config)
	if err != nil {
		self.fail = true
		return err
	}
	defer self.clean(func() {
		docker.ContainerRemove(container.ID)
	})()

	err = provider.AddIP(global_ip.String())
	if err != nil {
		self.fail = true
		return err
	}
	defer self.clean(func() {
		provider.RemoveIP(global_ip.String())
	})()

	err = provider.AddNat(global_ip.String(), local_ip.String())
	if err != nil {
		self.fail = true
		return err
	}
	defer self.clean(func() {
		provider.RemoveNat(global_ip.String(), local_ip.String())
	})()

	boxfile := boxfile.NewFromPath(util.BoxfileLocation())
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")
	planPayload := map[string]interface{}{"config": boxConfig.Parsed}
	jsonPayload, _ := json.Marshal(planPayload)

	// run plan hook TODO payload
	output, err := util.Exec(container.ID, "plan", string(jsonPayload))
	if err != nil {
		fmt.Println(output)
		self.fail = true
		return err
	}

	// save service in DB	
	service.ID = container.ID
	service.Name = self.config.Meta["name"]
	service.ExternalIP = global_ip.String()
	service.InternalIP = local_ip.String()

	err = json.Unmarshal([]byte(output), &service.Plan)
	if err != nil {
		self.fail = true
		return err
	}
	for i := 0; i < len(service.Plan.Users); i++ {
		service.Plan.Users[i].Password = util.RandomPassword()
	}

	// save the service 
	err = data.Put(util.AppName(), self.config.Meta["name"], service)
	if err != nil {
		self.fail = true
		return err
	}
	return nil
}