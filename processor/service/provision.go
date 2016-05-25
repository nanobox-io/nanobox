package service

import (
  "errors"
  "fmt"
	"os"

	"github.com/nanobox-io/nanobox-golang-stylish"

  "github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
  "github.com/nanobox-io/nanobox/util/data"
)

type serviceProvision struct {
  config  processor.ProcessConfig
  service svc
}

type svc struct {
  label   string
  name    string
  image   string
}

func init() {
  processor.Register("service_provision", serviceProvisionFunc)
}

func serviceProvisionFunc(config processor.ProcessConfig) (processor.Processor, error) {
  return serviceProvision{config: config}, nil
}

func (self serviceProvision) Results() processor.ProcessConfig {
  return self.config
}

func (self serviceProvision) Process() error {

  if err := self.validateMeta(); err != nil {
    return err
  }

  if err := self.buildService(); err != nil {
    return err
  }

  if err := self.processService(); err != nil {
    return err
  }

  return nil
}

// validateMeta validates that the provided metadata is supplied
func (self serviceProvision) validateMeta() error {

  if self.config.Meta["label"] == "" {
    return errors.New("missing service label")
  }

  if self.config.Meta["name"] == "" {
    return errors.New("missing service name")
  }

  if self.config.Meta["image"] == "" {
    return errors.New("missing service image")
  }

  return nil
}

// buildService builds the service struct
func (self *serviceProvision) buildService() error {

  self.service = svc{
    label:  self.config.Meta["label"],
    name:   self.config.Meta["name"],
    image:  self.config.Meta["image"],
  }

  return nil
}

// processService will either launch or start a service dependant on status
func (self serviceProvision) processService() error {
  active := self.isServiceActive()

  if active == true {
    return self.startService()
  } else {
    return self.launchService()
  }
}

// launchService will setup and configure a service
func (self serviceProvision) launchService() error {

	header := fmt.Sprintf("Provisioning %s...", self.service.label)
	fmt.Print(stylish.NestedBullet(header, self.config.DisplayLevel))

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		DisplayLevel: self.config.DisplayLevel + 1,
		Meta: map[string]string{
			"label": self.service.label,
			"name":  self.service.name,
			"image": self.service.image,
		},
	}

	// provision
	if err := processor.Run("service_setup", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_setup:", self.service.name), err)
		os.Exit(1)
	}

	// configure
	if err := processor.Run("service_configure", config); err != nil {
		fmt.Println(fmt.Sprintf("%s_setup:", self.service.name), err)
		os.Exit(1)
	}

	return nil
}

// startService will start a service
func (self serviceProvision) startService() error {

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"label": self.service.label,
			"name":  self.service.name,
		},
	}

	// start
	err := processor.Run("service_start", config)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s_start:", self.service.name), err)
		os.Exit(1)
	}

	return nil
}

// isServiceActive returns true if a service is already active
func (self serviceProvision) isServiceActive() bool {

  // service db entry
  service := models.Service{}

  // fetch the entry from the database, ignoring any errors as the service
  // might not exist yet
  data.Get(util.AppName(), self.service.name, &service)

	if service.State == "active" {
		return true
	}

	return false
}
