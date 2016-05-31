package processor

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type info struct {
	control ProcessControl
}

func init() {
	Register("dev_info", infoFunc)
}

func infoFunc(control ProcessControl) (Processor, error) {
	return info{control}, nil
}

func (self info) Results() ProcessControl {
	return self.control
}

func (self info) Process() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		fmt.Println("data keys:", err)
		lumber.Close()
		os.Exit(1)
	}

	for _, service := range services {
		if service != "builds" {
			svc := models.Service{}
			data.Get(util.AppName(), service, &svc)
			bytes, _ := json.MarshalIndent(svc, "", "  ")
			fmt.Printf("%s\n", bytes)
		}
	}

	envVars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &envVars)
	bytes, _ := json.MarshalIndent(envVars, "", "  ")
	fmt.Printf("%s\n", bytes)

	return nil
}
