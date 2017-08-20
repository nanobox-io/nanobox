package component

import (
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// Sync syncronizes an app's components with the boxfile config
func Sync(envModel *models.Env, appModel *models.App) error {
	display.OpenContext("Syncing data components")
	defer display.CloseContext()

	// clean any crufty that may have been left by a failed
	// sync
	if err := Clean(appModel); err != nil {
		return util.ErrorAppend(err, "failed to clean dirty components")
	}

	// purge delta components
	if err := purgeDeltaComponents(envModel, appModel); err != nil {
		return util.ErrorAppend(err, "failed to purge delta components")
	}

	// provision components
	if err := provisionComponents(envModel, appModel); err != nil {
		return util.ErrorAppend(err, "failed to provision components")
	}

	// update deployed boxfile
	appModel.DeployedBoxfile = envModel.BuiltBoxfile
	if err := appModel.Save(); err != nil {
		lumber.Error("component:Sync:models.App.Save()")
		return util.ErrorAppend(err, "failed to update deployed boxfile on app")
	}

	return nil
}

// purgeDeltaComponents purges components that have changed in the boxfile
func purgeDeltaComponents(envModel *models.Env, appModel *models.App) error {

	display.OpenContext("Removing old")
	defer display.CloseContext()

	// flag to determine if we display an up-to-date message
	upToDate := true

	// parse the boxfiles
	builtBoxfile := boxfile.New([]byte(envModel.BuiltBoxfile))
	deployedBoxfile := boxfile.New([]byte(appModel.DeployedBoxfile))

	components, err := models.AllComponentsByApp(appModel.ID)
	if err != nil {
		lumber.Error("component:purgeDeltaComponents:models.AllComponentsByApp(%s): %s", appModel.ID, err.Error())
		return util.ErrorAppend(err, "failed to load app components")
	}

	for _, component := range components {

		// ignore platform services
		if isPlatformUID(component.Name) {
			continue
		}

		// fetch the data nodes
		newNode := builtBoxfile.Node(component.Name)
		oldNode := deployedBoxfile.Node(component.Name)

		// skip if the new node is valid and they are the same
		if newNode.Valid && newNode.Equal(oldNode) {
			continue
		}

		upToDate = false

		// destroy the component
		if err := Destroy(appModel, component); err != nil {
			return util.ErrorAppend(err, "failed to destroy component")
		}
	}

	if upToDate {
		display.StartTask("Skipping (up-to-date)")
		display.StopTask()
	}

	return nil
}

// provisionComponents will provision components from the boxfile
func provisionComponents(envModel *models.Env, appModel *models.App) error {
	display.OpenContext("Launching new")
	defer display.CloseContext()

	// flag to determine if we display an up-to-date message
	upToDate := true

	// parse the boxfile
	builtBoxfile := boxfile.New([]byte(envModel.BuiltBoxfile))

	// grab all of the data nodes
	dataServices := builtBoxfile.Nodes("data")

	for _, name := range dataServices {
		// check to see if this component is already active
		componentModel, _ := models.FindComponentBySlug(appModel.ID, name)
		if componentModel.State == "active" {
			continue
		}

		upToDate = false

		componentModel.Name = name
		componentModel.Label = name
		componentModel.Image = builtBoxfile.Node(name).StringValue("image")

		// setup
		if err := Setup(appModel, componentModel); err != nil {
			// todo: if error `Error: No such image: image/postgresql` set code to USER, else, IMAGE
			return util.ErrorAppend(err, "failed to setup component (%s): %s", name, err.Error())
		}

	}

	if upToDate {
		display.StartTask("Skipping (up-to-date)")
		display.StopTask()
	}

	return nil
}

// isPlatform will return true if the uid matches a platform service
func isPlatformUID(uid string) bool {
	return uid == "portal" || uid == "hoarder" || uid == "mist" || uid == "logvac"
}
