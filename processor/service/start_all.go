package service

import (
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

	services, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		self.config.Meta["name"] = service
		err = processor.Run("service_start", self.config)
		if err != nil {
			return err
		}
	}
	return nil
}
