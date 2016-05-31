package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type devDestroy struct {
	control ProcessControl
}

func init() {
	Register("dev_destroy", devDestroyFunc)
}

func devDestroyFunc(control ProcessControl) (Processor, error) {
	return devDestroy{control}, nil
}

func (self devDestroy) Results() ProcessControl {
	return self.control
}

func (self devDestroy) Process() error {

	if err := Run("dev_setup", self.control); err != nil {
		return err
	}

	// remove all the services (platform/service/code)
	if err := self.removeServices(); err != nil {
		return err
	}

	// potentially destroy the platform
	if err := self.destroyPlatfrom(); err != nil {
		return err
	}

	return nil
}

// get all the services in the app
// and remove them
func (self devDestroy) removeServices() error {
	services, err := data.Keys(util.AppName())
	if err != nil {
		return fmt.Errorf("data keys: %s", err.Error())
	}
	self.control.Display(stylish.Bullet("Removing Services"))
	self.control.DisplayLevel++
	for _, service := range services {
		if service != "build" {
			// svc := models.Service{}
			// data.Get(util.AppName(), service, &svc)
			self.control.Meta["name"] = service
			err := Run("service_destroy", self.control)
			if err != nil {
				self.control.Display(stylish.Warning("one of the services did not uninstall:\n%s", err.Error()))
				// continue on to the next one. 
				// we should continue trying to remove services
			}
		}
	}	
	self.control.DisplayLevel--
	return nil
}

// if im the only app destroy the whole vm
func (self devDestroy) destroyPlatfrom() error {
	data.Delete("apps", util.AppName())
	keys, err := data.Keys("apps")
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		// if no other apps exist in container
		if err := Run("provider_destroy", self.control); err != nil {
			return err
		}
	}	
	return nil
}

