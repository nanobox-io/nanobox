package service

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStartAll struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("service_start_all", serviceStartAllFunc)
}

func serviceStartAllFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// make sure i was given a name and image
	return serviceStartAll{config: config}, nil
}

func (self serviceStartAll) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceStartAll) Process() error {

	if err := self.startServices(); err != nil {
		return err
	}

	return nil
}

// startServices starts all of the services saved in the database
func (self serviceStartAll) startServices() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := self.startService(service); err != nil {
			return err
		}
	}

	return nil
}

// startService starts a service
func (self serviceStartAll) startService(uid string) error {

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"label": 	uid,
			"name":  	uid,
		},
	}

	// provision
	if err := processor.Run("service_start", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_start:", uid), err)
		return err
	}

	return nil
}
