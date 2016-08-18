package code

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
)

// Setup ...
type Setup struct {
	App       models.App
	Component models.Component
	Name      string
	Image     string
	fail      bool
}

//
func (setup *Setup) clean(fn func()) func() {
	return func() {
		if setup.fail {
			fn()
		}
	}
}

//
func (setup *Setup) Run() error {
	// quit early if the service was found in the database as well as docker
	if setup.componentExists() {
		return nil
	}

	display.OpenContext("setting up %s", setup.Name)
	defer display.CloseContext()

	//
	if err := setup.getLocalIP(); err != nil {
		lumber.Error("code:Setup:setup.getLocalIP(): %s", err.Error())
		setup.fail = true
		return err
	}
	defer setup.clean(setup.returnLocalIP)()

	//
	if err := setup.getGlobalIP(); err != nil {
		lumber.Error("code:Setup:setup.getGlobalIP(): %s", err.Error())
		setup.fail = true
		return err
	}
	defer setup.clean(setup.returnGlobalIP)()

	// pull the docker image
	streamer := display.NewStreamer("info")
	dockerPercent := &display.DockerPercentDisplay{Output: streamer, Prefix: setup.Image}
	if _, err := docker.ImagePull(setup.Image, dockerPercent); err != nil {
		lumber.Error("code:Setup:docker.ImagePull(%s, nil): %s", setup.Image, err.Error())
		setup.fail = true
		return err
	}

	//
	if err := setup.createContainer(); err != nil {
		return err
	}
	defer setup.clean(setup.removeContainer)()

	//
	if err := setup.addIPToProvider(); err != nil {
		return err
	}
	defer setup.clean(setup.removeIPFromProvider)()

	// save the component
	if err := setup.Component.Save(); err != nil {
		lumber.Error("code:Setup:Component.Save(): %s", err.Error())
		setup.fail = true
		return err
	}

	return nil
}

// componentExists ...
func (setup *Setup) componentExists() bool {
	var databaseErr error
	setup.Component, databaseErr = models.FindComponentBySlug(setup.App.ID, setup.Name)

	_, containerErr := docker.GetContainer(setup.Component.ID)

	return databaseErr == nil && containerErr == nil
}

// getLocalIP ...
func (setup *Setup) getLocalIP() error {

	// if the service already has an ip
	if setup.Component.InternalIP != "" {
		return nil
	}
	ip, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("code:Setup:dhcp.ReserveLocal(): %s", err.Error())
		setup.fail = true
		return err
	}
	setup.Component.InternalIP = ip.String()

	return nil
}

// getGlobalIP ...
func (setup *Setup) getGlobalIP() error {
	if setup.Component.ExternalIP != "" {
		// if the service already has an ip
		return nil
	}
	ip, err := dhcp.ReserveGlobal()
	if err != nil {
		lumber.Error("code:Setup:dhcp.ReserveGlobal(): %s", err.Error())
		setup.fail = true
		return err
	}
	setup.Component.ExternalIP = ip.String()

	return nil
}

// returnLocalIP ...
func (setup *Setup) returnLocalIP() {
	dhcp.ReturnIP(net.ParseIP(setup.Component.InternalIP))
}

// returnGlobalIP ...
func (setup *Setup) returnGlobalIP() {
	dhcp.ReturnIP(net.ParseIP(setup.Component.ExternalIP))
}

// addIPToProvider ...
func (setup *Setup) addIPToProvider() error {
	if setup.Component.InternalIP == "" || setup.Component.ExternalIP == "" {
		return fmt.Errorf("im missing an ip to bind to the provider")
	}

	if err := provider.AddIP(setup.Component.ExternalIP); err != nil {
		lumber.Error("code:Setup:addIPToProvider:provider.AddIP(%s): %s", setup.Component.ExternalIP, err.Error())
		setup.fail = true
		return err
	}

	if err := provider.AddNat(setup.Component.ExternalIP, setup.Component.InternalIP); err != nil {
		lumber.Error("code:Setup:addIPToProvider:provider.AddNat(%s, %s): %s", setup.Component.ExternalIP, setup.Component.InternalIP, err.Error())
		setup.fail = true
		return err
	}
	return nil
}

// removeIPFromProvider ...
func (setup *Setup) removeIPFromProvider() {
	provider.RemoveNat(setup.Component.ExternalIP, setup.Component.InternalIP)
	provider.RemoveIP(setup.Component.ExternalIP)
}

// createContainer ...
func (setup *Setup) createContainer() error {
	display.StartTask("creating container")

	// configure the container
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s_%s", setup.App.ID, setup.Name),
		Image:   setup.Image,
		Network: "virt",
		IP:      setup.Component.InternalIP,
	}

	// create docker container
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Setup:createContainer:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		setup.fail = true
		return err
	}
	setup.Component.AppID = setup.App.ID
	setup.Component.ID = container.ID
	setup.Component.Name = setup.Name
	setup.Component.Type = "code"

	display.StopTask()
	return nil
}

// removeContainer ...
func (setup *Setup) removeContainer() {
	// catch the error here and display it but dont error out
	docker.ContainerRemove(setup.Component.ID)
}
