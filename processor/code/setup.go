package code

import (
	"errors"
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
)

type codeSetup struct {
	control  processor.ProcessControl
	service models.Service
	fail    bool
}

var missingImageOrName = errors.New("missing image or name")

func init() {
	processor.Register("code_setup", codeSetupFunc)
}

func codeSetupFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeSetup{control: control}, nil
}

func (self *codeSetup) clean(fn func()) func() {
	return func() {
		if self.fail {
			fn()
		}
	}
}

func (self codeSetup) Results() processor.ProcessControl {
	return self.control
}

func (self *codeSetup) Process() error {
	// make sure i was given a name and image
	if self.control.Meta["name"] == "" || self.control.Meta["image"] == "" {
		return missingImageOrName
	}

	if self.serviceExists() {
		// quit early if the service was found in the database as well as docker
		return nil
	}

	if err := self.getLocalIP(); err != nil {
		return err
	}
	defer self.clean(self.returnLocalIP)()

	if err := self.getGlobalIP(); err != nil {
		return err
	}
	defer self.clean(self.returnGlobalIP)()

	// pull the docker image
	prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(self.control.DisplayLevel), self.control.Meta["image"])

	if _, err := docker.ImagePull(self.control.Meta["image"], &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
		return err
	}

	if err := self.createContainer(); err != nil {
		return err
	}
	defer self.clean(self.removeContainer)()

	if err := self.addIPToProvider(); err != nil {
		return err
	}
	defer self.clean(self.removeIPFromProvider)()

	// save the service
	if err := data.Put(util.AppName(), self.control.Meta["name"], self.service); err != nil {
		self.fail = true
		lumber.Error("insert data: ", err)
		return err
	}
	return nil
}

func (self *codeSetup) serviceExists() bool {
	service := models.Service{}
	databaseErr := data.Get(util.AppName(), self.control.Meta["name"], &service)
	if databaseErr == nil {
		// set the service i found so i dont re allocate ips
		self.service = service
	}
	_, containerErr := docker.GetContainer(service.ID)
	return databaseErr == nil && containerErr == nil
}

func (self *codeSetup) getLocalIP() error {
	if self.service.InternalIP != "" {
		// if the service already has an ip
		return nil
	}
	ip, err := ip_control.ReserveLocal()
	if err != nil {
		self.fail = true
		return err
	}
	self.service.InternalIP = ip.String()
	return nil
}

func (self *codeSetup) getGlobalIP() error {
	if self.service.ExternalIP != "" {
		// if the service already has an ip
		return nil
	}
	ip, err := ip_control.ReserveGlobal()
	if err != nil {
		self.fail = true
		return err
	}
	self.service.ExternalIP = ip.String()
	return nil
}

func (self *codeSetup) returnLocalIP() {
	ip_control.ReturnIP(net.ParseIP(self.service.InternalIP))
}

func (self *codeSetup) returnGlobalIP() {
	ip_control.ReturnIP(net.ParseIP(self.service.ExternalIP))
}

func (self *codeSetup) addIPToProvider() error {
	if self.service.InternalIP == "" || self.service.ExternalIP == "" {
		return fmt.Errorf("im missing an ip to bind to the provider")
	}

	if err := provider.AddIP(self.service.ExternalIP); err != nil {
		self.fail = true
		return err
	}

	if err := provider.AddNat(self.service.ExternalIP, self.service.InternalIP); err != nil {
		self.fail = true
		return err
	}
	return nil
}

func (self *codeSetup) removeIPFromProvider() {
	provider.RemoveNat(self.service.ExternalIP, self.service.InternalIP)
	provider.RemoveIP(self.service.ExternalIP)
}

func (self *codeSetup) createContainer() error {
	// configure the container
	fmt.Println("-> building container", self.control.Meta["name"])
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox-%s-%s", util.AppName(), self.control.Meta["name"]),
		Image:   self.control.Meta["image"],
		Network: "virt",
		IP:      self.service.InternalIP,
	}

	// create docker container
	container, err := docker.CreateContainer(config)
	if err != nil {
		self.fail = true
		lumber.Error("container: ", err)
		return err
	}
	self.service.ID = container.ID
	self.service.Name = self.control.Meta["name"]
	self.service.Type = "code"
	return nil
}

func (self *codeSetup) removeContainer() {
	// catch the error here and display it but dont error out
	docker.ContainerRemove(self.service.ID)
}
