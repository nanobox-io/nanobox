package service

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceStopAll struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("service_stop_all", serviceStopAllFunc)
}

func serviceStopAllFunc(config processor.ProcessConfig) (processor.Processor, error) {
	return serviceStopAll{config: config}, nil
}

func (self serviceStopAll) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceStopAll) Process() error {

	if err := self.stopServices(); err != nil {
		return err
	}

	return nil
}

// stopServices stops all of the services saved in the database
func (self serviceStopAll) stopServices() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := self.stopService(service); err != nil {
			return err
		}
	}

	return nil
}

// stopService stops a service
func (self serviceStopAll) stopService(uid string) error {

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"label": 	uid,
			"name":  	uid,
		},
	}

	// provision
	if err := processor.Run("service_stop", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_stop:", uid), err)
		return err
	}

	return nil
}
