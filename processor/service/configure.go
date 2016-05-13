package service

import (
	"encoding/json"
	"fmt"
	"errors"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dockerexec"
)

type serviceConfigure struct {
	config processor.ProcessConfig
	service models.Service
}

type member struct {
	LocalIP string `json:"local_ip"`
	UID     string `json:"uid"`
	Role    string `json:"role"`
}

type component struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	ID   string `json:"id"`
}

type configPayload struct {
	LogvacHost string                 `json:"logvac_host"`
	Platform   string                 `json:"platform"`
	Config     map[string]interface{} `json:"config"`
	Member     member                 `json:"member"`
	Component  component              `json:"component"`
	Users      []models.User          `json:"users"`
}

type startPayload struct {
	Config map[string]interface{} `json:"config"`
}

func init() {
	processor.Register("service_configure", serviceConfigureFunc)
}

func (self serviceConfigure) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceConfigure) Process() error {

	if err := self.validateImage(); err != nil {
		return err
	}

	if err := self.loadService(); err != nil {
		return err
	}

	// short-circuit if the service is already started
	if self.service.Started {
		return nil
	}

	if err := self.runUpdate(); err != nil {
		return nil
	}

	if err := self.runConfigure(); err != nil {
		return nil
	}

	if err := self.runStart(); err != nil {
		return nil
	}

	if err := self.persistService(); err != nil {
		return nil
	}

	return nil
}

func serviceConfigureFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return serviceConfigure{config: config}, nil
}

func (self serviceConfigure) configurePayload() string {
	me := models.Service{}
	data.Get(util.AppName(), self.config.Meta["name"], &me)

	logvac := models.Service{}
	data.Get(util.AppName(), "logvac", &logvac)

	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")

	pload := configPayload{
		LogvacHost: logvac.InternalIP,
		Platform:   "local",
		Config:     boxConfig.Parsed,
		Member: member{
			LocalIP: me.InternalIP,
			UID:     "1",
			Role:    "primary",
		},
		Component: component{
			Name: "whydoesthismatter",
			UID:  self.config.Meta["name"],
			ID:   me.ID,
		},
		Users: me.Plan.Users,
	}
	if pload.Users == nil {
		pload.Users = []models.User{}
	}
	switch self.config.Meta["name"] {
	case "portal", "logvac", "hoarder", "mist":
		pload.Config["token"] = "123"
	}
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(j)
}

func (self serviceConfigure) startPayload() string {
	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")

	pload := startPayload{boxConfig.Parsed}
	switch self.config.Meta["name"] {
	case "portal", "logvac", "hoarder", "mist":
		pload.Config["token"] = "123"
	}
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(j)
}

// validateImage validates that the image is provided
func (self *serviceConfigure) validateImage() error {
	// make sure i was given a name and image
	if self.config.Meta["name"] == "" {
		return errors.New("missing image or name")
	}

	return nil
}

// loadService loads the service from the database
func (self *serviceConfigure) loadService() error {
	// get the service from the database
	err := data.Get(util.AppName(), self.config.Meta["name"], &self.service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// runUpdate will run the update hook in the container
func (self *serviceConfigure) runUpdate() error {
	// run update
	fmt.Print(stylish.NestedBullet("Updating...", self.config.DisplayLevel))

	cmd := dockerexec.Command(self.service.ID, "update", "{}")
	if err := cmd.Run(); err != nil {
		fmt.Println(cmd.Output())
		return err
	}

	return nil
}

// runConfigure will run the configure hook in the container
func (self *serviceConfigure) runConfigure() error {
	// run update
	fmt.Print(stylish.NestedBullet("Configuring services...", self.config.DisplayLevel))

	cmd := dockerexec.Command(self.service.ID, "configure", "{}")
	if err := cmd.Run(); err != nil {
		fmt.Println(cmd.Output())
		return err
	}

	return nil
}

// runStart will run the configure hook in the container
func (self *serviceConfigure) runStart() error {
	// run update
	fmt.Print(stylish.NestedBullet("Starting services...", self.config.DisplayLevel))

	cmd := dockerexec.Command(self.service.ID, "start", "{}")
	if err := cmd.Run(); err != nil {
		fmt.Println(cmd.Output())
		return err
	}

	return nil
}

// persistService saves the service entry to the database
func (self *serviceConfigure) persistService() error {
	// set started to true for further idempotency
	self.service.Started = true
	err := data.Put(util.AppName(), self.config.Meta["name"], &self.service)
	if err != nil {
		return err
	}

	return nil
}
