package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

//
func Setup(appModel *models.App, componentModel *models.Component) error {
	display.OpenContext("setting up %s", componentModel.Name)
	defer display.CloseContext()

	// generate the missing component data
	if err := componentModel.Generate(appModel, "code"); err != nil {
		lumber.Error("component:Setup:models.Component:Generate(%s, code): %s", appModel.ID, componentModel.Name, err.Error())
		return fmt.Errorf("failed to generate component data: %s", err.Error())
	}

	// short-circuit if this component is already setup
	if componentModel.State != "initialized" {
		return nil
	}

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		Prefix: componentModel.Image,
	}

	// pull the component image
	if _, err := docker.ImagePull(componentModel.Image, dockerPercent); err != nil {
		lumber.Error("component:Setup:docker.ImagePull(%s, nil): %s", componentModel.Image, err.Error())
		return fmt.Errorf("failed to pull docker image (%s): %s", componentModel.Image, err.Error())
	}

	//

	if err := reserveIps(componentModel); err != nil {
		lumber.Error("code:Setup:setup.getLocalIP(): %s", err.Error())
		return err
	}


	// create docker container
	config := generate_container.ComponentConfig(componentModel)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Setup:createContainer:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return err
	}

	// save the component
	componentModel.ID = container.ID
	if err := componentModel.Save(); err != nil {
		lumber.Error("code:Setup:Component.Save(): %s", err.Error())
		return err
	}

	//
	if err := addIPToProvider(componentModel); err != nil {
		return err
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

// createContainer ...
func createContainer(componentModel *models.Component) error {
	display.StartTask("creating container")
	

	display.StopTask()
	
	return nil
}

// addIPToProvider ...
func addIPToProvider(componentModel *models.Component) error {

	if err := provider.AddIP(componentModel.ExternalIP); err != nil {
		lumber.Error("code:Setup:addIPToProvider:provider.AddIP(%s): %s", componentModel.ExternalIP, err.Error())
		return err
	}

	if err := provider.AddNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("code:Setup:addIPToProvider:provider.AddNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		return err
	}
	return nil
}
