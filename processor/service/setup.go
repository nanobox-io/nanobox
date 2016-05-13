package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"io"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/dockerexec"
)

type serviceSetup struct {
	config 		processor.ProcessConfig
	service 	models.Service
	local_ip 	net.IP
	global_ip net.IP
	container dockType.ContainerJSON
	plan			string
	fail   		bool
}

func init() {
	processor.Register("service_setup", serviceSetupFunc)
}

func serviceSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &serviceSetup{config: config}, nil
}

func (self *serviceSetup) clean(fn func()) func() {
	return func() {
		if self.fail {
			fn()
		}
	}
}

func (self serviceSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceSetup) Process() error {

	if err := self.validateImage(); err != nil {
		return err
	}

	if err := self.loadService(); err != nil {
		return err
	}

	// quit early if the service was found to be created already
	if self.service.ID != "" {
		return nil
	}

	if err := self.downloadImage(); err != nil {
		return err
	}

	if err := self.reserveIps(); err != nil {
		return err
	}

	if err := self.launchContainer(); err != nil {
		return err
	}

	if err := self.attachNetwork(); err != nil {
		return err
	}

	if err := self.planService(); err != nil {
		return err
	}

	if err := self.persistService(); err != nil {
		return err
	}

	return nil
}

// validateImage ensures we were given a name and image
func (self *serviceSetup) validateImage() error {
	if self.config.Meta["name"] == "" || self.config.Meta["image"] == "" {
		return errors.New("missing image or name")
	}
	return nil
}

// loadService fetches the service from the database
func (self *serviceSetup) loadService() error {
	return data.Get(util.AppName(), self.config.Meta["name"], &self.service)
}

// downloadImage downloads the docker image
func (self *serviceSetup) downloadImage() error {
	label := "Downloading docker image " + self.config.Meta["image"]
	fmt.Print(stylish.NestedProcessStart(label, self.config.DisplayLevel))

	// Create a pipe to for the JSONMessagesStream to read from
	pr, pw := io.Pipe()
	prefix := stylish.GenerateNestedPrefix(self.config.DisplayLevel + 1)
  go print.DisplayJSONMessagesStream(pr, os.Stdout, os.Stdout.Fd(), true, prefix, nil)
	if _, err := docker.ImagePull(self.config.Meta["image"], pw); err != nil {
		return err
	}
  fmt.Print(stylish.ProcessEnd())

	return nil
}

// reserveIps reserves a global and local ip for the container
func (self *serviceSetup) reserveIps() error {

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

	// assign back to the state
	self.local_ip = local_ip
	self.global_ip = global_ip

	return nil
}

// launchContainer launches and starts a docker container
func (self *serviceSetup) launchContainer() error {
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("%s-%s", util.AppName(), self.config.Meta["name"]),
		Image:   self.config.Meta["image"],
		Network: "virt",
		IP:      self.local_ip.String(),
	}

	fmt.Print(stylish.NestedProcessStart("Starting docker container...", self.config.DisplayLevel))
	container, err := docker.CreateContainer(config)
	if err != nil {
		self.fail = true
		lumber.Error("container: ", err)
		return err
	}
	defer self.clean(func() {
		docker.ContainerRemove(container.ID)
	})()

	self.container = container

	return nil
}

// attachNetwork attaches the IP addresses to the container
func (self *serviceSetup) attachNetwork() error {
	label := "Add container to host network..."
	fmt.Print(stylish.NestedProcessStart(label, self.config.DisplayLevel))

	err := provider.AddIP(self.global_ip.String())
	if err != nil {
		self.fail = true
		lumber.Error("addip: ", err)
		return err
	}
	defer self.clean(func() {
		provider.RemoveIP(self.global_ip.String())
	})()

	err = provider.AddNat(self.global_ip.String(), self.local_ip.String())
	if err != nil {
		self.fail = true
		lumber.Error("addnat: ", err)
		return err
	}
	defer self.clean(func() {
		provider.RemoveNat(self.global_ip.String(), self.local_ip.String())
	})()

	return nil
}

// planService runs the plan hook
func (self *serviceSetup) planService() error {
	fmt.Print(stylish.NestedProcessStart("Gathering service requirements...", self.config.DisplayLevel))

	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")
	planPayload := map[string]interface{}{"config": boxConfig.Parsed}
	jsonPayload, _ := json.Marshal(planPayload)


	cmd := dockerexec.Command(self.container.ID, "plan", string(jsonPayload))
	if err := cmd.Run(); err != nil {
		fmt.Println(cmd.Output())
		self.fail = true
		lumber.Error("plan: ", err)
		return err
	}

	self.plan = cmd.Output()

	return nil
}

// persistService saves the service in the database
func (self *serviceSetup) persistService() error {
	// save service in DB
	self.service.ID = self.container.ID
	self.service.Name = self.config.Meta["name"]
	self.service.ExternalIP = self.global_ip.String()
	self.service.InternalIP = self.local_ip.String()

	err := json.Unmarshal([]byte(self.plan), &self.service.Plan)
	if err != nil {
		self.fail = true
		return err
	}
	for i := 0; i < len(self.service.Plan.Users); i++ {
		self.service.Plan.Users[i].Password = util.RandomString(10)
	}

	// save the service
	err = data.Put(util.AppName(), self.config.Meta["name"], &self.service)
	if err != nil {
		self.fail = true
		lumber.Error("insert data: ", err)
		return err
	}

	return nil
}
