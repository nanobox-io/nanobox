package processor

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type info struct {
	config ProcessConfig
}

func init() {
	Register("info", infoFunc)
}

func infoFunc(config ProcessConfig) (Processor, error) {
	return info{config}, nil
}

func (self info) Results() ProcessConfig {
	return self.config
}

func (self info) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

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
			fmt.Printf("%+v\n", svc)
		}
	}

	return nil
}
