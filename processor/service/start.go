package service

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStart struct {
	config processor.ProcessConfig
}

// {
//   "logvac_host": "127.0.0.1",
//   "platform": "local",
//   "config": {
//     "token": "123"
//   },
//   "member": {
//     "local_ip": "192.168.0.2",
//     "uid": "1",
//     "role": "primary"
//   },
//   "component": {
//     "name": "willy-walrus",
//     "uid": "logvac1",
//     "id": "9097d0a7-7e02-4be5-bce1-3d7cb1189488"
//   },
//   "users": [

//   ]
// }
type configPayload struct {
	LogvacHost string `json:"logvac_host"`
	Platform   string `json:"platform"`
	Config     map[string]interface{} `json:"config"`
	Member     struct {
		LocalIP string `json:"local_ip"`
		UID     string `json:"uid"`
		Role    string `json:"role"`
	} `json:"member"`
	Component struct {
		Name string `json:"name"`
		UID  string `json:"uid"`
		ID   string `json:"id"`
	}
	Users []models.User `json:"users"`
}

func init() {
	processor.Register("serivce_start", serviceStartFunc)
	docker.Initialize("env")
}

func serviceStartFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return serviceStart{config: config}, nil
}

func (self serviceStart) configurePayload() string {
	me := models.Service{}
	data.Get(util.AppName(), self.config.Meta["name"], &me)

	logvac := models.Service{}
	data.Get(util.AppName(), "logvac", &logvac)

	boxfile := boxfile.NewFromPath(util.BoxfileLocation())
	boxConfig := boxfile.Node(self.config.Meta["name"]).Node("config")


	pload := configPayload{
		LogvacHost: logvac.InternalIP,
		Config: boxConfig.Parsed,
		Users: me.Plan.Users,
	}
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(j)

}

func (self serviceStart) startPayload() string {
	return ""
}

func (self serviceStart) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceStart) Process() error {
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

	// run configure command TODO PAYLOAD
	output, err := util.Exec(service.ID, "configure", self.configurePayload())
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