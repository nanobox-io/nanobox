package component

import (
	"regexp"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Sync ...
type Sync struct {
	Env             models.Env
	App             models.App
	builtBoxfile    boxfile.Boxfile
	deployedBoxfile boxfile.Boxfile
}

//
func (sync *Sync) Run() error {
	display.OpenContext("syncronizing boxfile components")
	defer display.CloseContext()

	if err := sync.loadBuiltBoxfile(); err != nil {
		return err
	}

	if err := sync.loadDeployedBoxfile(); err != nil {
		return err
	}

	if err := sync.purgeDeltaComponents(); err != nil {
		return err
	}

	if err := sync.provisionComponents(); err != nil {
		return err
	}

	if err := sync.updateDeployedBoxfile(); err != nil {
		return err
	}

	return nil
}

// loadBuiltBoxfile loads the new build boxfile from the database
func (sync *Sync) loadBuiltBoxfile() error {
	sync.builtBoxfile = boxfile.New([]byte(sync.Env.BuiltBoxfile))

	return nil
}

// loadDeployedBoxfile loads the old boxfile from the database
func (sync *Sync) loadDeployedBoxfile() error {
	sync.deployedBoxfile = boxfile.New([]byte(sync.App.DeployedBoxfile))

	return nil
}

// update that we have deployed
func (sync *Sync) updateDeployedBoxfile() error {
	sync.App.DeployedBoxfile = sync.Env.BuiltBoxfile
	return sync.App.Save()
}

// purgeDeltaComponents will purge the services that were removed from the boxfile
func (sync *Sync) purgeDeltaComponents() error {
	display.StartTask("perging delta components")

	components, err := models.AllComponentsByApp(sync.App.ID)
	if err != nil {
		display.ErrorTask()
		return err
	}

	for _, component := range components {

		// ignore platform services
		if isPlatformUID(component.Name) {
			continue
		}

		// fetch the nodes
		newNode := sync.builtBoxfile.Node(component.Name)
		oldNode := sync.deployedBoxfile.Node(component.Name)

		// skip if the new node is valid and they are the same
		if newNode.Valid && newNode.Equal(oldNode) {
			continue
		}

		if err := sync.purgeComponent(component); err != nil {
			display.ErrorTask()
			return err
		}

	}

	display.StopTask()
	return nil
}

// purgeComponent will purge a service from the nanobox
func (sync *Sync) purgeComponent(component models.Component) error {
	destroy := Destroy{sync.App, component}
	return destroy.Run()
}

// provisionServices will provision services that are defined in the boxfile
// but not running on nanobox
func (sync *Sync) provisionComponents() error {
	display.StartTask("building new/updated components")

	// grab all of the data nodes
	dataServices := sync.builtBoxfile.Nodes("data")

	for _, name := range dataServices {
		image := sync.builtBoxfile.Node(name).StringValue("image")

		if image == "" {
			serviceType := regexp.MustCompile(`.+\.`).ReplaceAllString(name, "")
			image = "nanobox/" + serviceType
		}

		// check to see if this component is already active
		comp, _ := models.FindComponentBySlug(sync.App.ID, name)
		if comp.State == ACTIVE {
			continue
		}

		setup := &Setup{
			App:   sync.App,
			Image: image,
			Name:  name,
		}

		// setup the service
		if err := setup.Run(); err != nil {
			display.ErrorTask()
			return err
		}

		// each component has the potential to update the app
		sync.App = setup.App

		configure := Configure{
			App:       sync.App,
			Component: setup.Component,
		}

		// and configure it
		if err := configure.Run(); err != nil {
			display.ErrorTask()
			return err
		}

	}

	display.StopTask()
	return nil
}

// isPlatform will return true if the uid matches a platform service
func isPlatformUID(uid string) bool {
	return uid == PORTAL || uid == HOARDER || uid == MIST || uid == LOGVAC
}
