package processor

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type devDestroy struct {
	config ProcessConfig
}

func init() {
	Register("dev_destroy", devDestroyFunc)
}

func devDestroyFunc(config ProcessConfig) (Processor, error) {
	return devDestroy{config}, nil
}

func (self devDestroy) Results() ProcessConfig {
	return self.config
}

func (self devDestroy) Process() error {

	// if im the only app dont even worry about any of the service
	// clean up just destroy the whole vm
	data.Delete("apps", util.AppName())
	keys, err := data.Keys("apps")
	if err != nil {
		fmt.Println("get apps data failure:", err)
		lumber.Close()
		os.Exit(1)
	}
	if len(keys) == 0 {
		// if no other apps exist in container
		err := Run("provider_destroy", self.config)
		if err != nil {
			fmt.Println("provider_setup:", err)
			lumber.Close()
			os.Exit(1)
		}
		return nil
	}

	// get all the services in the app
	// and remove them
	services, err := data.Keys(util.AppName())
	if err != nil {
		fmt.Println("data keys:", err)
		lumber.Close()
		os.Exit(1)
	}

	for _, service := range services {
		if service != "build" {
			svc := models.Service{}
			data.Get(util.AppName(), service, &svc)
			self.config.Meta["name"] = service
			err := Run("service_destroy", self.config)
			if err != nil {
				fmt.Println("remove service failure:", err)
				lumber.Close()
				os.Exit(1)
			}
		}
	}
	return nil
}
