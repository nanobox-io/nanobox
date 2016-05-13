package service

import (
	"errors"

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
	// make sure i was given a name and image
	if config.Meta["name"] == "" {
		return nil, errors.New("missing image or name")
	}

	return serviceStopAll{config: config}, nil
}

func (self serviceStopAll) Results() processor.ProcessConfig {
	return self.config
}

func (self serviceStopAll) Process() error {

	services, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, service := range services {
		self.config.Meta["name"] = service
		err = processor.Run("service_stop", self.config)
		if err != nil {
			return err
		}
	}
	return nil
}
