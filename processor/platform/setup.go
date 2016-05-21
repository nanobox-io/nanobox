package platform

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
)

type platformSetup struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("platform_setup", platformSetupFunc)
}

func platformSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return platformSetup{config}, nil
}

func (self platformSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self platformSetup) Process() error {

	if err := self.provisionServices(); err != nil {
		return err
	}

	return nil
}

// provisionServices will provision all the platform services
func (self platformSetup) provisionServices() error {
	for _, service := range PlatformServices {
		if err := self.provisionService(service); err != nil {
			return err
		}
	}

	return nil
}

// provisionService will provision an individual service
func (self platformSetup) provisionService(service PlatformService) error {

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"label": 	service.label,
			"name":  	service.name,
			"image":	service.image,
		},
	}

	// provision
	if err := processor.Run("service_provision", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_provision:", service.name), err)
		return err
	}

	return nil
}
