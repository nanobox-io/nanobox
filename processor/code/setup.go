package code

import (
	"errors"
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type codeSetup struct {
	config processor.ProcessConfig
	fail bool
}

var missingName = errors.New("missing name")

func init() {
	processor.Register("code_setup", codeSetupFunc)
}

func codeSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeSetup{config: config}, nil
}

func (self *codeSetup) clean(fn func()) func() {
	return func() {
		if self.fail {
			fn()
		}
	}
}

func (self codeSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self *codeSetup) Process() error {
	// make sure i was given a name and image
	if self.config.Meta["name"] == ""
		return missingName
	}

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &service)

	// create docker container
	if err == nil {
		// quit early if the service was found to be created already
		return nil
	}

	image := self.config.Meta["image"]
	if image == "" {
		iamge = "nanobox/code"
	}

	_, err = docker.ImagePull(image)
	if err != nil {
		return err
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
		Image: image,
 		Network: "virt",
 		IP: local_ip.String(),
	}

	container, err := docker.CreateContainer(config)
	if err != nil {
		self.fail = true
		lumber.Error("container: ", err)
		return err
	}
	defer self.clean(func() {
		docker.ContainerRemove(container.ID)
	})()

	err = provider.AddIP(global_ip.String())
	if err != nil {
		self.fail = true
		lumber.Error("addip: ", err)
		return err
	}
	defer self.clean(func() {
		provider.RemoveIP(global_ip.String())
	})()

	err = provider.AddNat(global_ip.String(), local_ip.String())
	if err != nil {
		self.fail = true
		lumber.Error("addnat: ", err)
		return err
	}
	defer self.clean(func() {
		provider.RemoveNat(global_ip.String(), local_ip.String())
	})()

	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")

	// run plan hook TODO payload
	output, err := util.Exec(container.ID, "plan", string(jsonPayload))
	if err != nil {
		fmt.Println(output)
		self.fail = true
		lumber.Error("plan: ", err)
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
		lumber.Error("insert data: ", err)
		return err
	}
	lumber.Debug("worked")
	return nil
}