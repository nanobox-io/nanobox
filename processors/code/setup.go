package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/code"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/util/provider"
)

//
func Setup(appModel *models.App, componentModel *models.Component, warehouseConfig WarehouseConfig) error {
	// generate the missing component data
	if err := componentModel.Generate(appModel, "code"); err != nil {
		lumber.Error("component:Setup:models.Component:Generate(%s, code): %s", appModel.ID, componentModel.Name, err.Error())
		return fmt.Errorf("failed to generate component data: %s", err.Error())
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
		if _, err := docker.ImagePull(componentModel.Image, dockerPercent); err != nil {
			lumber.Error("component:Setup:docker.ImagePull(%s, nil): %s", componentModel.Image, err.Error())
			display.ErrorTask()
			return fmt.Errorf("failed to pull docker image (%s): %s", componentModel.Image, err.Error())
		}
		display.StopTask()

	}

	display.StartTask("Starting docker container")
	if err := reserveIps(componentModel); err != nil {
		display.ErrorTask()
		lumber.Error("code:Setup:setup.getLocalIP(): %s", err.Error())
		return err
	}

	// create docker container
	config := container_generator.ComponentConfig(componentModel)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Setup:createContainer:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// save the component
	componentModel.ID = container.ID
	if err := componentModel.Save(); err != nil {
		lumber.Error("code:Setup:Component.Save(): %s", err.Error())
		return err
	}

	// attach container to the host network
	if err := attachNetwork(componentModel); err != nil {
		return fmt.Errorf("failed to attach container to host network: %s", err.Error())
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
		return fmt.Errorf("failed to configure code: %s", err.Error())
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
		lumber.Error("code:Configure:Component.Save(): %s", err.Error())
		return fmt.Errorf("unable to save component model: %s", err.Error())
	}

	return nil
}

//  ...
func reserveIps(componentModel *models.Component) error {
	if componentModel.InternalIP == "" {
		localIP, err := dhcp.ReserveLocal()
		if err != nil {
			lumber.Error("code:Setup:dhcp.ReserveLocal(): %s", err.Error())
			return err
		}
		componentModel.InternalIP = localIP.String()
	}

	if componentModel.ExternalIP == "" {
		ip, err := dhcp.ReserveGlobal()
		if err != nil {
			lumber.Error("code:Setup:dhcp.ReserveGlobal(): %s", err.Error())
			return err
		}
		componentModel.ExternalIP = ip.String()
	}

	return nil
}

// attachNetwork attaches the component to the host network
func attachNetwork(componentModel *models.Component) error {
	if err := provider.AddIP(componentModel.ExternalIP); err != nil {
		display.ErrorTask()
		lumber.Error("code:Setup:addIPToProvider:provider.AddIP(%s): %s", componentModel.ExternalIP, err.Error())
		return err
	}

	if err := provider.AddNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("code:Setup:addIPToProvider:provider.AddNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		display.ErrorTask()
		return err
	}

	return nil
}
