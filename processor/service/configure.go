package service

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceConfigure struct {
	config processor.ProcessConfig
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

func (self serviceConfigure) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceConfigure) Process() error {
	// make sure i was given a name and image
	if self.config.Meta["name"] == "" {
		return missingImageOrName
	}

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	if service.Started {
		return nil
	}

	// run update
	output, err := util.Exec(service.ID, "update", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	// run configure command TODO PAYLOAD
	output, err = util.Exec(service.ID, "configure", self.configurePayload())
	if err != nil {
		fmt.Println(output)
		return err
	}

	// run start command TODO PAYLOAD
	output, err = util.Exec(service.ID, "start", self.startPayload())
	if err != nil {
		fmt.Println(output)
		return err
	}

	service.Started = true
	err = data.Put(util.AppName(), self.config.Meta["name"], service)
	if err != nil {
		return err
	}

	return nil
}
