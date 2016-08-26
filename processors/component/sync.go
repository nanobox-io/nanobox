package component

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
)

// Sync syncronizes an app's components with the boxfile config
func Sync(envModel *models.Env, appModel *models.App) error {

	// purge delta components
	if err := purgeDeltaComponents(envModel, appModel); err != nil {
		return fmt.Errorf("failed to purge delta components: %s", err.Error())
	}

	// provision components
	if err := provisionComponents(envModel, appModel); err != nil {
		return fmt.Errorf("failed to provision components: %s", err.Error())
	}

	// update deployed boxfile
	appModel.DeployedBoxfile = envModel.BuiltBoxfile
	if err := appModel.Save(); err != nil {
		lumber.Error("component:Sync:models.App.Save(): %s", err.Error())
		return fmt.Errorf("failed to update deployed boxfile on app: %s", err.Error())
	}

	return nil
}

// purgeDeltaComponents purges components that have changed in the boxfile
func purgeDeltaComponents(envModel *models.Env, appModel *models.App) error {
	// parse the boxfiles
	builtBoxfile := boxfile.New([]byte(envModel.BuiltBoxfile))
	deployedBoxfile := boxfile.New([]byte(appModel.DeployedBoxfile))

	components, err := models.AllComponentsByApp(appModel.ID)
	if err != nil {
		lumber.Error("component:purgeDeltaComponents:models.AllComponentsByApp(%s): %s", appModel.ID, err.Error())
		return fmt.Errorf("failed to load app components: %s", err.Error())
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

		// destroy the component
		if err := Destroy(appModel, component); err != nil {
			return fmt.Errorf("failed to destroy component: %s", err.Error())
		}
	}

	return nil
}

// provisionComponents will provision components from the boxfile
func provisionComponents(envModel *models.Env, appModel *models.App) error {
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

		componentModel.Name = name
		componentModel.Image = builtBoxfile.Node(name).StringValue("image")

		// setup
		if err := Setup(appModel, componentModel); err != nil {
			return fmt.Errorf("failed to setup component (%s): %s", name, err.Error())
		}

		// configure
		if err := Configure(appModel, componentModel); err != nil {
			return fmt.Errorf("failed to configure component: %s", err.Error())
		}
	}

	return nil
}

// isPlatform will return true if the uid matches a platform service
func isPlatformUID(uid string) bool {
	return uid == "portal" || uid == "hoarder" || uid == "mist" || uid == "logvac"
}
