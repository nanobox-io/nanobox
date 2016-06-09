package service

import (
	"encoding/json"
	"errors"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceConfigure struct {
	control  	processor.ProcessControl
	service 	models.Service
	boxfile 	models.Boxfile
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

type startUpdatePayload struct {
	Config map[string]interface{} `json:"config"`
}

func init() {
	processor.Register("service_configure", serviceConfigureFunc)
}

func (self serviceConfigure) Results() processor.ProcessControl {
	return self.control
}

func (self serviceConfigure) Process() error {

	if err := self.validateMeta(); err != nil {
		return err
	}

	if err := self.loadService(); err != nil {
		return err
	}

	// short-circuit if the service has already progressed past this point
	if self.service.State != "planned" {
		return nil
	}

	if err := self.loadBoxfile(); err != nil {
		return err
	}

	if err := self.runUpdate(); err != nil {
		return err
	}

	if err := self.runConfigure(); err != nil {
		return err
	}

	if err := self.runStart(); err != nil {
		return err
	}

	if err := self.persistService(); err != nil {
		return err
	}

	return nil
}

func serviceConfigureFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return serviceConfigure{control: control}, nil
}

func (self serviceConfigure) configurePayload() string {
	me := models.Service{}
	data.Get(util.AppName(), self.control.Meta["name"], &me)

	logvac := models.Service{}
	data.Get(util.AppName(), "logvac", &logvac)

	box := boxfile.New([]byte(self.boxfile.Data))
	boxConfig := box.Node(self.control.Meta["name"]).Node("config")

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
			UID:  self.control.Meta["name"],
			ID:   me.ID,
		},
		Users: me.Plan.Users,
	}
	if pload.Users == nil {
		pload.Users = []models.User{}
	}
	switch self.control.Meta["name"] {
	case "portal", "logvac", "hoarder", "mist":
		pload.Config["token"] = "123"
	}
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(j)
}

func (self serviceConfigure) startUpdatePayload() string {
	boxfile := boxfile.New([]byte(self.control.Meta["boxfile"]))
	boxConfig := boxfile.Node(self.control.Meta["name"]).Node("config")

	pload := startUpdatePayload{boxConfig.Parsed}
	switch self.control.Meta["name"] {
	case "portal", "logvac", "hoarder", "mist":
		pload.Config["token"] = "123"
	}
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(j)
}

// validateMeta validates that the image is provided
func (self *serviceConfigure) validateMeta() error {
	// make sure i was given a name and image
	if self.control.Meta["name"] == "" {
		return errors.New("missing name")
	}

	return nil
}

// loadService loads the service from the database
func (self *serviceConfigure) loadService() error {
	// get the service from the database
	err := data.Get(util.AppName(), self.control.Meta["name"], &self.service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	return nil
}

// loadBoxfile loads the new build boxfile from the database
func (self *serviceConfigure) loadBoxfile() error {

	if err := data.Get(util.AppName()+"_meta", "build_boxfile", &self.boxfile); err != nil {
		return err
	}

	return nil
}

// runUpdate will run the update hook in the container
func (self *serviceConfigure) runUpdate() error {
	self.control.Info(stylish.SubBullet("Updating services..."))

	// run update
	_, err := util.Exec(self.service.ID, "update", self.startUpdatePayload(), nil)
	return err
}

// runConfigure will run the configure hook in the container
func (self *serviceConfigure) runConfigure() error {
	// run configure
	self.control.Info(stylish.SubBullet("Configuring services..."))

	_, err := util.Exec(self.service.ID, "configure", self.configurePayload(), nil)
	return err
}

// runStart will run the configure hook in the container
func (self *serviceConfigure) runStart() error {
	// run update
	self.control.Info(stylish.SubBullet("Starting services..."))
	_, err := util.Exec(self.service.ID, "start", self.startUpdatePayload(), nil)
	return err
}

// persistService saves the service entry to the database
func (self *serviceConfigure) persistService() error {
	self.service.State = "active"
	err := data.Put(util.AppName(), self.control.Meta["name"], &self.service)
	if err != nil {
		return err
	}

	return nil
}
