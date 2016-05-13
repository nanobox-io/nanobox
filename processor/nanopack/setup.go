package nanopack

import (
	"fmt"
	"os"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
)

type nanopackSetup struct {
	config processor.ProcessConfig
}

type nanoService struct {
	label 	string
	name 		string
	image 	string
}

var (
	services = []nanoService{
		nanoService{
			label: 	"Load Balancer",
			name: 	"portal",
			image: 	"nanobox/portal",
		},
		nanoService{
			label: 	"Realtime Message Bus",
			name: 	"mist",
			image: 	"nanobox/mist",
		},
		nanoService{
			label: 	"Logger",
			name: 	"logvac",
			image: 	"nanobox/logvac",
		},
		nanoService{
			label: 	"Warehouse",
			name: 	"hoarder",
			image: 	"nanobox/hoarder",
		},
	}
)

func init() {
	processor.Register("nanopack_setup", nanopackSetupFunc)
}

func nanopackSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return nanopackSetup{config}, nil
}

func (self nanopackSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self nanopackSetup) Process() error {

	// let's short-circuit if the platform is already up and running
	if running := isPlatformRunning(); running == true {
		return nil
	}

	label := "Provisioning Platform Services..."
	fmt.Print(stylish.NestedBullet(label, self.config.DisplayLevel))

	for _, service := range services {
		if err := self.processService(&service); err != nil {
			return err
		}
	}

	return nil
}

// processService will process an individual service
func (self nanopackSetup) processService(service *nanoService) error {

	created := isPlatformServiceCreated(service.name)

	if created == true {
		return self.startService(service)
	} else {
		return self.launchService(service)
	}
}

// launchService will setup and configure a service
func (self nanopackSetup) launchService(service *nanoService) error {
	header := fmt.Sprintf("Launching %s...", service.label)
	fmt.Print(stylish.NestedBullet(header, self.config.DisplayLevel + 1))

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		DisplayLevel: self.config.DisplayLevel + 2,
		Meta: map[string]string{
			"label": service.label,
			"name":  service.name,
			"image": service.image,
		},
	}

	// provision
	err := processor.Run("service_setup", config)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s_setup:", service.name), err)
		os.Exit(1)
	}

	// configure
	err = processor.Run("service_configure", config)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s_setup:", service.name), err)
		os.Exit(1)
	}

	return nil
}

// startService will start a service
func (self nanopackSetup) startService(service *nanoService) error {

	// short-circuit if this service is already running
	if running := isPlatformServiceRunning(service.name); running == true {
		return nil
	}

	header := fmt.Sprintf("Booting %s...", service.label)
	fmt.Print(stylish.NestedBullet(header, self.config.DisplayLevel + 1))

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		DisplayLevel: self.config.DisplayLevel + 2,
		Meta: map[string]string{
			"label": service.label,
			"name":  service.name,
		},
	}

	// start
	err := processor.Run("service_start", config)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s_start:", service.name), err)
		os.Exit(1)
	}

	return nil
}

// isPlatformRunning returns true if all the platform services are running
func isPlatformRunning() bool {

	for _, service := range services {
		if running := isPlatformServiceRunning(service.name); running == false {
			return false
		}
	}

	return true
}

// isPlatformServiceCreated returns true if a service is already created
func isPlatformServiceCreated(service string) bool {
	name := fmt.Sprintf("%s-%s", util.AppName(), service)

	_, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	if err != nil {
		return false
	}

	return true
}

// isPlatformServiceRunning returns true if a service is already running
func isPlatformServiceRunning(service string) bool {
	name := fmt.Sprintf("%s-%s", util.AppName(), service)

	container, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	if err != nil {
		return false
	}

	// return true if the container is running
	if container.State.Status == "running" {
		return true
	}

	return false
}
