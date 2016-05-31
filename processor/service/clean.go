package service

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceClean struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("service_clean", serviceCleanFunc)
}

func serviceCleanFunc(control processor.ProcessControl) (processor.Processor, error) {
	return serviceClean{control: control}, nil
}

func (self serviceClean) Results() processor.ProcessControl {
	return self.control
}

func (self serviceClean) Process() error {

	if err := self.cleanServices(); err != nil {
		self.control.Display(stylish.Warning("there has been an error cleaning a service:\n%s", err.Error()))
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
	self.control.Display(stylish.Bullet(fmt.Sprintf("Cleaning %s...", uid)))

	config := processor.ProcessControl{
		DevMode:      self.control.DevMode,
		Verbose:      self.control.Verbose,
		DisplayLevel: self.control.DisplayLevel + 1,
		Meta: map[string]string{
			"name": uid,
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

	return !containerExists(uid)
}

// containerExists will check to see if a docker container exists on the provider
func containerExists(uid string) bool {
	_, err := docker.GetContainer(fmt.Sprintf("nanobox-%s-%s", util.AppName(), uid))
	return err == nil
}
