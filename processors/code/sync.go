package code

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Sync is used by the deploy process to syncronize the code parts
func Sync(appModel *models.App, warehouseConfig WarehouseConfig) error {
	display.OpenContext("starting code components")
	defer display.CloseContext()

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// iterate over the code nodes and build containers for each of them
	for _, componentModel := range codeComponentModels(appModel) {

		// run the code setup process with the new config
		err := Setup(appModel, componentModel)
		if err != nil {
			return fmt.Errorf("failed to setup code (%s): %s\n", componentModel.Name, err.Error())
		}

		// configure this code container
		err = Configure(appModel, componentModel, warehouseConfig)
		if err != nil {
			return fmt.Errorf("failed to configure code (%s): %s\n", componentModel.Name, err.Error())
		}
	}

	return nil
}

// setBoxfile ...
func codeComponentModels(appModel *models.App) []*models.Component {

	componentModels := []*models.Component{}

	// look in the boxfile for code nodes and generate a stub component
	box := boxfile.New([]byte(appModel.DeployedBoxfile))
	for _, componentName := range box.Nodes("code") {
		image := box.Node(componentName).StringValue("image")
		if image == "" {
			image = "nanobox/code:v1"
		}

		componentModel := &models.Component{
			Name: componentName,
			Image: image,
		}

		componentModels = append(componentModels, componentModel)
	}

	return componentModels
}
