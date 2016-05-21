package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"io"
	"os"

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

type cleanFunc func() error

type serviceSetup struct {
	config 			processor.ProcessConfig
	service 		models.Service
	local_ip 		net.IP
	global_ip 	net.IP
	container 	dockType.ContainerJSON
	plan				string
	fail   			bool
	cleanFuncs	[]cleanFunc
}

func init() {
	processor.Register("service_setup", serviceSetupFunc)
}

func serviceSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &serviceSetup{
		config: config,
		cleanFuncs:	make([]cleanFunc, 0),
	}, nil
}

// clean will iterate through the cleanup functions that were registered and
// call them one-by-one
func (self *serviceSetup) clean() error {
	// short-circuit if we haven't failed
	if self.fail == false {
		return nil
	}

	// iterate through the cleanup functions that were registered and call them
	for _, cleanF := range self.cleanFuncs {
		if err := cleanF(); err != nil {
			return err
		}
	}

	return nil
}

func (self serviceSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceSetup) Process() error {

	// call the cleanup function to ensure we don't leave any bad state
	defer self.clean()

	if err := self.validateMeta(); err != nil {
		self.fail = true
		return err
	}

	if err := self.loadService(); err != nil {
		self.fail = true
		return err
	}

	// short-circuit if the service has already progressed past this point
	if self.service.State != "initialized" {
		return nil
	}

	if err := self.downloadImage(); err != nil {
		self.fail = true
		return err
	}

	if err := self.reserveIps(); err != nil {
		self.fail = true
		return err
	}

	if err := self.launchContainer(); err != nil {
		self.fail = true
		return err
	}

	if err := self.attachNetwork(); err != nil {
		self.fail = true
		return err
	}

	if err := self.planService(); err != nil {
		self.fail = true
		return err
	}

	if err := self.persistService(); err != nil {
		self.fail = true
		return err
	}

	return nil
}

// validateMeta ensures we were given a name and image
func (self *serviceSetup) validateMeta() error {
	if self.config.Meta["name"] == "" || self.config.Meta["image"] == "" {
		return errors.New("missing image or name")
	}
	return nil
}

// loadService fetches the service from the database
func (self *serviceSetup) loadService() error {
	// the service really shouldn't exist yet, so let's not return the error if it fails
	data.Get(util.AppName(), self.config.Meta["name"], &self.service)

	// set the default state
	if self.service.State == "" {
		self.service.State = "initialized"
	}

	return nil
}

// downloadImage downloads the docker image
func (self *serviceSetup) downloadImage() error {
	label := "Pulling image " + self.config.Meta["image"]
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

	self.cleanFuncs = append(self.cleanFuncs, func() error {
		return ip_control.ReturnIP(local_ip)
	})

	global_ip, err := ip_control.ReserveGlobal()
	if err != nil {
		return err
	}

	self.cleanFuncs = append(self.cleanFuncs, func() error {
		return ip_control.ReturnIP(global_ip)
	})

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

	fmt.Print(stylish.NestedBullet("Starting container...", self.config.DisplayLevel))
	container, err := docker.CreateContainer(config)
	if err != nil {
		return err
	}

	self.cleanFuncs = append(self.cleanFuncs, func() error {
		return docker.ContainerRemove(container.ID)
	})

	self.container = container

	return nil
}

// attachNetwork attaches the IP addresses to the container
func (self *serviceSetup) attachNetwork() error {
	label := "Bridging container to host network..."
	fmt.Print(stylish.NestedBullet(label, self.config.DisplayLevel))

	err := provider.AddIP(self.global_ip.String())
	if err != nil {
		return err
	}

	self.cleanFuncs = append(self.cleanFuncs, func() error {
		return provider.RemoveIP(self.global_ip.String())
	})

	err = provider.AddNat(self.global_ip.String(), self.local_ip.String())
	if err != nil {
		return err
	}

	self.cleanFuncs = append(self.cleanFuncs, func() error {
		return provider.RemoveNat(self.global_ip.String(), self.local_ip.String())
	})

	return nil
}

// planService runs the plan hook
func (self *serviceSetup) planService() error {
	fmt.Print(stylish.NestedBullet("Gathering service requirements...", self.config.DisplayLevel))

	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")
	planPayload := map[string]interface{}{"config": boxConfig.Parsed}
	jsonPayload, _ := json.Marshal(planPayload)


	cmd := dockerexec.Command(self.container.ID, "plan", string(jsonPayload))
	if err := cmd.Run(); err != nil {
		fmt.Println(cmd.Output())
		return err
	}

	self.plan = cmd.Output()

	return nil
}

// persistService saves the service in the database
func (self *serviceSetup) persistService() error {
	// save service in DB
	self.service.ID         = self.container.ID
	self.service.Name       = self.config.Meta["name"]
	self.service.ExternalIP = self.global_ip.String()
	self.service.InternalIP = self.local_ip.String()
	self.service.State      = "planned"

	err := json.Unmarshal([]byte(self.plan), &self.service.Plan)
	if err != nil {
		return err
	}
	for i := 0; i < len(self.service.Plan.Users); i++ {
		self.service.Plan.Users[i].Password = util.RandomString(10)
	}

	// save the service
	err = data.Put(util.AppName(), self.config.Meta["name"], &self.service)
	if err != nil {
		return err
	}

	return nil
}
