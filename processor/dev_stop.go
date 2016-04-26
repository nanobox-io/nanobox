package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type devStop struct {
	config ProcessConfig
}

func init() {
	Register("dev_stop", devStopFunc)
}

func devStopFunc(config ProcessConfig) (Processor, error) {
	return devStop{config}, nil
}

func (self devStop) Results() ProcessConfig {
	return self.config
}

func (self devStop) Process() error {
	app := models.App{}
	data.Get("apps", util.AppName(), app)

	if app.UsageCount == 0 {
		// remove all code containers
		boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))

		for _, codeName := range boxfile.Nodes("code") {

			self.config.Meta["name"] = codeName
			// dont catch errors because i dont
			// want destroy to fail
			err := Run("code_destroy", self.config)
			if err != nil {
				fmt.Printf("code_destroy (%s): %s\n", codeName, err.Error())
			}

		}
		// stop all data containers
		Run("service_stop_all", self.config)
	}

	//if nothing is running stop the provider
	keys, _ := data.Keys("apps")
	for _, key := range keys {
		data.Get("apps", key, &app)
		if app.UsageCount != 0 {
			// i found an app that is doing something
			// so dont shut down the provider
			return nil
		}
	}

	// this part only gets executed if we cant find any apps
	// doing anything in the loop above
	return Run("provider_stop", self.config)
	// return nil
}
