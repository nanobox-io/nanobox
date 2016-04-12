package code

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type codeStart struct {
	config processor.ProcessConfig
}

type startPayload struct {
	Config map[string]interface{} `json:"config"`
}

func init() {
	processor.Register("code_start", codeStartFunc)
}

func codeStartFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return codeStart{config: config}, nil
}

func (self codeStart) startPayload() string {
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

func (self codeStart) Results() processor.ProcessConfig {
	return self.config
}

func (self codeStart) Process() error {
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

	// run fetch build command
	output, err = util.Exec(service.ID, "fetch", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	// run start command
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
