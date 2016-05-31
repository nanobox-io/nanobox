package service

import (
	"fmt"
	"regexp"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type serviceSync struct {
	control     processor.ProcessControl
	fail       bool
	newBoxfile models.Boxfile
	oldBoxfile models.Boxfile
}

type service struct {
	label string
	name  string
	image string
}

func init() {
	processor.Register("service_sync", serviceSyncFunc)
}

func serviceSyncFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &serviceSync{control: control}, nil
}

func (self serviceSync) Results() processor.ProcessControl {
	return self.control
}

func (self *serviceSync) Process() error {

	if err := self.loadNewBoxfile(); err != nil {
		return err
	}

	if err := self.loadOldBoxfile(); err != nil {
		return err
	}

	if err := self.purgeDeltaServices(); err != nil {
		return err
	}

	if err := self.provisionDataServices(); err != nil {
		return err
	}

	if err := self.replaceOldBoxfile(); err != nil {
		return err
	}

	return nil
}

// loadNewBoxfile loads the new build boxfile from the database
func (self *serviceSync) loadNewBoxfile() error {

	if err := data.Get(util.AppName()+"_meta", "build_boxfile", &self.newBoxfile); err != nil {
		return err
	}

	return nil
}

// loadOldBoxfile loads the old boxfile from the database
func (self *serviceSync) loadOldBoxfile() error {

	// we don't care about the error here because this could be the first build
	data.Get(util.AppName()+"_meta", "old_build_boxfile", &self.oldBoxfile)

	return nil
}

// replaceOldBoxfile replaces the old boxfile in the database with the new boxfile
func (self *serviceSync) replaceOldBoxfile() error {

	if err := data.Put(util.AppName()+"_meta", "old_build_boxfile", &self.newBoxfile); err != nil {
		return err
	}

	return nil
}

// purgeDeltaServices will purge the services that were removed from the boxfile
func (self *serviceSync) purgeDeltaServices() error {

	// convert the data into boxfile library helpers
	oldBoxfile := boxfile.New(self.oldBoxfile.Data)
	newBoxfile := boxfile.New(self.newBoxfile.Data)

	// fetch the services
	uids, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	for _, uid := range uids {

		// ignore platform services
		if isPlatformUID(uid) {
			continue
		}

		// fetch the nodes
		newNode := newBoxfile.Node(uid)
		oldNode := oldBoxfile.Node(uid)

		// skip if the nodes are the same
		if newNode.Equal(oldNode) {
			continue
		}

		if err := self.purgeService(uid); err != nil {
			return err
		}

	}

	return nil
}

// purgeService will purge a service from the nanobox
func (self *serviceSync) purgeService(uid string) error {
	service := processor.ProcessControl{
		DevMode: self.control.DevMode,
		Verbose: self.control.Verbose,
		Meta: map[string]string{
			"name": uid,
		},
	}

	if err := processor.Run("service_remove", service); err != nil {
		fmt.Printf("service_remove (%s): %s\n", uid, err.Error())
		return err
	}

	return nil
}

// provisionServices will provision services that are defined in the boxfile
// but not running on nanobox
func (self *serviceSync) provisionDataServices() error {

	// convert the data into boxfile library helpers
	newBoxfile := boxfile.New(self.newBoxfile.Data)

	// grab all of the data nodes
	dataServices := newBoxfile.Nodes("data")

	for _, uid := range dataServices {
		image := newBoxfile.Node(uid).StringValue("image")

		if image == "" {
			serviceType := regexp.MustCompile(`.+\.`).ReplaceAllString(uid, "")
			image = "nanobox/" + serviceType
		}

		config := processor.ProcessControl{
			DevMode:      self.control.DevMode,
			Verbose:      self.control.Verbose,
			DisplayLevel: self.control.DisplayLevel + 1,
			Meta: map[string]string{
				"name":  uid,
				"image": image,
			},
		}

		// setup the service
		if err := processor.Run("service_setup", config); err != nil {
			return err
		}

		// and configure it
		if err := processor.Run("service_configure", config); err != nil {
			return err
		}

	}

	return nil
}

// isPlatform will return true if the uid matches a platform service
func isPlatformUID(uid string) bool {
	return uid == "portal" || uid == "hoarder" || uid == "mist" || uid == "logvac"
}
