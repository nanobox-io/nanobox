package platform

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
)

type platformStop struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("platform_stop", platformStopFunc)
}

func platformStopFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return platformStop{config}, nil
}

func (self platformStop) Results() processor.ProcessConfig {
	return self.config
}

func (self platformStop) Process() error {

	if err := self.stopServices(); err != nil {
		return err
	}

	return nil
}

// stopServices will stop all the platform services
func (self *platformStop) stopServices() error {
	for _, service := range PlatformServices {
		if err := self.stopService(service); err != nil {
			return err
		}
	}

	return nil
}

// stopService will stop an individual service
func (self *platformStop) stopService(service PlatformService) error {

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"label": 	service.label,
			"name":  	service.name,
		},
	}

	// stop
	if err := processor.Run("service_stop", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_stop:", service.name), err)
		return err
	}

	return nil
}
