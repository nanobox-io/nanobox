package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceSync struct {
	config processor.ProcessConfig
	fail   bool
}

func init() {
	processor.Register("service_sync", serviceSyncFunc)
}

func serviceSyncFunc(config processor.ProcessConfig) (processor.Processor, error) {
	if config.Meta["boxfile"] == "" {
		return nil, errors.New("missing boxfile")
	}

	return &serviceSync{config: config}, nil
}

func (self serviceSync) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceSync) Process() error {
	// populate new boxfile
	box := boxfile.New([]byte(self.config.Meta["boxfile"]))

	// get the previous boxfile
	oldBoxData := models.Boxfile{}
	data.Get(util.AppName()+"_meta", "oldBoxfile", &oldBoxData)
	oldBoxfile := boxfile.New(oldBoxData.Data)

	// remove all the services no longer in the boxfile
	// or a change has happened to the boxfile node
	keys, err := data.Keys(util.AppName())
	if err != nil {
		fmt.Println(err)
	}
	for _, key := range keys {
		// if the boxfile doesnt have a node for the service
		// or if the old and new boxfiles dont match for this service
		if key != "portal" &&
			key != "hoarder" &&
			key != "mist" &&
			key != "logvac" &&
			(!box.Node(key).Valid ||
				box.Node(key).Equal(oldBoxfile.Node(key))) {
			service := processor.ProcessConfig{
				DevMode: self.config.DevMode,
				Verbose: self.config.Verbose,
				Meta: map[string]string{
					"name":    key,
					"boxfile": self.config.Meta["boxfile"],
				},
			}
			err := processor.Run("service_remove", service)
			if err != nil {
				fmt.Printf("service_remove (%s): %s\n", key, err.Error())
				os.Exit(1)
			}
		}
	}

	// add any missing services
	for _, serviceName := range box.Nodes("data") {
		image := box.Node(serviceName).StringValue("image")
		if image == "" {
			serviceType := regexp.MustCompile(`.+\.`).ReplaceAllString(serviceName, "")
			image = "nanobox/" + serviceType
		}
		service := processor.ProcessConfig{
			DevMode: self.config.DevMode,
			Verbose: self.config.Verbose,
			Meta: map[string]string{
				"name":    serviceName,
				"image":   image,
				"boxfile": self.config.Meta["boxfile"],
			},
		}
		err := processor.Run("service_setup", service)
		if err != nil {
			fmt.Printf("service_setup (%s): %s\n", serviceName, err.Error())
			os.Exit(1)
		}

		err = processor.Run("service_configure", service)
		if err != nil {
			fmt.Printf("service_setup (%s): %s\n", serviceName, err.Error())
			os.Exit(1)
		}
	}

	// set the new box file as the old for next time we sync
	return data.Put(util.AppName()+"_meta", "oldBoxfile", box)
}
