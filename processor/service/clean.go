package service

import (
  "fmt"

  "github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/golang-docker-client"

  "github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
  "github.com/nanobox-io/nanobox/util/data"
)

type serviceClean struct {
  config processor.ProcessConfig
}

func init() {
  processor.Register("service_clean", serviceCleanFunc)
}

func serviceCleanFunc(config processor.ProcessConfig) (processor.Processor, error) {
  return serviceClean{config: config}, nil
}

func (self serviceClean) Results() processor.ProcessConfig {
  return self.config
}

func (self serviceClean) Process() error {

  if err := self.cleanServices(); err != nil {
    return nil
  }

  return nil
}

// cleanServices will iterate through each registered service and
// clean them if they were left in a bad state
func (self serviceClean) cleanServices() error {

  uids, err := data.Keys(util.AppName())
  if err != nil {
    return err
  }

  for _, uid := range uids {
    if err := self.cleanService(uid); err != nil {
      return err
    }
  }

  return nil
}

// cleanService will clean a service if it was left in a bad state
func (self serviceClean) cleanService(uid string) error {

  if dirty := isServiceDirty(uid); dirty == true {
    return self.removeService(uid)
  }

  return nil
}

// removeService will remove a service from nanobox
func (self serviceClean) removeService(uid string) error {
  header := fmt.Sprintf("Cleaning %s...", uid)
	fmt.Print(stylish.NestedBullet(header, self.config.DisplayLevel))

  config := processor.ProcessConfig{
    DevMode: self.config.DevMode,
    Verbose: self.config.Verbose,
    DisplayLevel: self.config.DisplayLevel + 1,
    Meta: map[string]string{
      "name":  uid,
    },
  }

  err := processor.Run("service_remove", config)
  if err != nil {
    fmt.Println(fmt.Sprintf("%s_remove:", uid), err)
    return err
  }

  return nil
}

// isServiceDirty will return true if the service is not active and available
func isServiceDirty(uid string) bool {
  // service db entry
  service := models.Service{}

  // fetch the entry from the database
  if err := data.Get(util.AppName(), uid, &service); err != nil {
    return true
  }

  // short-circuit if this service never made it to active
  if service.State != "active" {
    return true
  }

  if exists := containerExists(uid); exists != true {
    return true
  }

  return false
}

// containerExists will check to see if a docker container exists on the provider
func containerExists(uid string) bool {
  name := fmt.Sprintf("%s-%s", util.AppName(), uid)

  if _, err := docker.GetContainer(name); err == nil {
    return true
  }

  return false
}
