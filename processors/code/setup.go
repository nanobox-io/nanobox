package code

import (
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/code"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

//
func Setup(appModel *models.App, componentModel *models.Component, warehouseConfig WarehouseConfig) error {
	// generate the missing component data
	if err := componentModel.Generate(appModel, "code"); err != nil {
		lumber.Error("component:Setup:models.Component:Generate(%s, code): %s", appModel.ID, componentModel.Name, err.Error())
		return util.ErrorAppend(err, "failed to generate component data")
	}

	// short-circuit if this component is already setup
	if componentModel.State != "initialized" {
		return nil
	}

	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: componentModel.Image,
	}

	if !docker.ImageExists(componentModel.Image) {
		// pull the component image
		display.StartTask("Pulling %s image", componentModel.Image)
		imagePull := func() error {
			_, err := docker.ImagePull(componentModel.Image, dockerPercent)
			return err
		}
		if err := util.Retry(imagePull, 5, time.Second); err != nil {
			lumber.Error("component:Setup:docker.ImagePull(%s, nil): %s", componentModel.Image, err.Error())
			display.ErrorTask()
			return util.ErrorAppend(err, "failed to pull docker image (%s)", componentModel.Image)
		}
		display.StopTask()

	}

	display.StartTask("Starting docker container")
	if err := reserveIps(componentModel); err != nil {
		display.ErrorTask()
		lumber.Error("code:Setup:setup.getLocalIP()")
		return err
	}

	// create docker container
	config := container_generator.ComponentConfig(componentModel)
	// remove any container that may have been created with this name befor
	// this can happen if the process is killed after the 
	// container was created but before our db model was saved
	docker.ContainerRemove(config.Name)
	
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Setup:createContainer:docker.CreateContainer(%+v)", config)
		display.ErrorTask()
		return util.ErrorAppend(err, "unable to create container")
	}
	display.StopTask()

	// save the component
	componentModel.ID = container.ID
	if err := componentModel.Save(); err != nil {
		lumber.Error("code:Setup:Component.Save()")
		return err
	}

	lumber.Prefix("code:Setup")
	defer lumber.Prefix("")

	// run fetch build command
	fetchPayload := hook_generator.FetchPayload(componentModel, warehouseConfig.WarehouseURL)

	display.StartTask("Fetching build from warehouse")
	if _, err := hookit.DebugExec(componentModel.ID, "fetch", fetchPayload, "info"); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// run configure command
	payload := hook_generator.ConfigurePayload(appModel, componentModel)

	//
	display.StartTask("Starting services")
	if _, err := hookit.DebugExec(componentModel.ID, "configure", payload, "info"); err != nil {
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to configure code")
	}

	// run start command
	if _, err := hookit.DebugExec(componentModel.ID, "start", payload, "info"); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	//
	componentModel.State = ACTIVE
	if err := componentModel.Save(); err != nil {
		lumber.Error("code:Configure:Component.Save()")
		return util.ErrorAppend(err, "unable to save component model")
	}

	return nil
}

//  ...
func reserveIps(componentModel *models.Component) error {
	if componentModel.IPAddr() == "" {
		localIP, err := dhcp.ReserveLocal()
		if err != nil {
			lumber.Error("code:Setup:dhcp.ReserveLocal()")
			return err
		}
		componentModel.IP = localIP.String()
	}

	return nil
}
