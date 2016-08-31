package dev

import (
	"fmt"
	"net"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	build_generator "github.com/nanobox-io/nanobox/generators/hooks/build"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	process_provider "github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/watch"
)

//
func Console(appModel *models.App, devRun bool) error {

	// run the share init which gives access to docker
	if err := process_provider.Init(); err != nil {
		return err
	}

	// generate a container config
	config := container_generator.DevConfig(appModel)

	defer teardown(appModel, config)

	// setup the docker container one is needed
	if err := setup(appModel, config); err != nil {
		return err
	}

	envModel, _ := models.FindEnvByID(appModel.EnvID)

	go watch.Watch(envModel.Directory)

	// if run then start the run commands
	// and log but do not continue to the regular console
	if devRun {
		return Run(appModel)
	}

	//
	if err := printMOTD(appModel); err != nil {
		return err
	}

	//
	if err := runConsole(appModel); err != nil {
		return err
	}

	return nil
}

// setup ...
func setup(appModel *models.App, config docker.ContainerConfig) error {

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// let anyone else know we're using the dev container
	counter.Increment(appModel.ID)

	//
	if !isDevRunning() {

		//
		if err := downloadImage(config.Image); err != nil {
			return err
		}

		container, err := docker.CreateContainer(config)
		if err != nil {
			return err
		}

		//
		if err := attachNetwork(appModel, config); err != nil {
			return err
		}

		lumber.Prefix("dev:Console")
		defer lumber.Prefix("")

		// TODO: The nil in this call needs to be replaced with something from the new display
		userPayload := build_generator.UserPayload()
		if _, err := hookit.Exec(container.ID, "user", userPayload, "debug"); err != nil {
			return err
		}

		// TODO: The nil in this call needs to be replaced with something from the new display
		if _, err := hookit.Exec(container.ID, "dev", build_generator.DevPayload(appModel), "debug"); err != nil {
			return err
		}
	} else {

		// if im not creating one i need to release the ip that was reserved in the config
		dhcp.ReturnIP(net.ParseIP(config.IP))
	}

	return nil
}

// teardown ...
func teardown(appModel *models.App, config docker.ContainerConfig) error {

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	counter.Decrement(appModel.ID)

	//
	if devIsUnused() {

		//
		if err := removeContainer(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}

		//
		if err := detachNetwork(appModel, config); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}

		//
		if err := dhcp.ReturnIP(net.ParseIP(config.IP)); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}
	}

	return nil
}

// downloadImage downloads the dev docker image
func downloadImage(image string) error {
	if !docker.ImageExists(image) {

		streamer := display.NewStreamer("info")
		dockerPercent := &display.DockerPercentDisplay{Output: streamer, Prefix: image}

		if _, err := docker.ImagePull(image, dockerPercent); err != nil {
			return err
		}
	}

	return nil
}

// removeContainer will lookup the dev container and remove it
func removeContainer() error {

	// grab the container info
	container, err := docker.GetContainer(container_generator.DevName())
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}

	if err := docker.ContainerRemove(container.ID); err != nil {
		// but if the container exists and we cant remove it for some other reason
		// we need to report that error
		return err
	}

	return nil
}

// runConsole will establish a console within the dev container
func runConsole(appModel *models.App) error {

	// create a dummy component using the appname
	component := &models.Component{
		ID: "nanobox_" + appModel.ID,
	}
	// for tyler: I dont like forcing someone into zsh..
	// your chosen shell is a very personal
	consoleConfig := env.ConsoleConfig{
		Cwd: cwd(appModel),
	}

	return env.Console(component, consoleConfig)
}

// cwd sets the cwd from the boxfile or provides a sensible default
func cwd(appModel *models.App) string {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	// parse the boxfile data

	if boxfile.Node("dev").StringValue("cwd") != "" {
		return boxfile.Node("dev").StringValue("cwd")
	}

	return "/app"
}

// printMOTD prints the motd with information for the user to connect
func printMOTD(appModel *models.App) error {

	// fetch the dev ip
	devIP := appModel.GlobalIPs["env"]

	os.Stderr.WriteString(fmt.Sprintf(`
                                   **
                                ********
                             ***************
                          *********************
                            *****************
                          ::    *********    ::
                             ::    ***    ::
                           ++   :::   :::   ++
                              ++   :::   ++
                                 ++   ++
                                    +
                    _  _ ____ _  _ ____ ___  ____ _  _
                    |\ | |__| |\ | |  | |__) |  |  \/
                    | \| |  | | \| |__| |__) |__| _/\_

--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm
+ Changes to your code in either the vm or workstation will be mirrored
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------

`, devIP))

	return nil
}

// attachNetwork attaches the container to the host network
func attachNetwork(appModel *models.App, config docker.ContainerConfig) error {

	// fetch the devIP
	devIP := appModel.GlobalIPs["env"]

	//
	if err := provider.AddIP(devIP); err != nil {
		return fmt.Errorf("provider:add_ip: %s", err.Error())
	}

	//
	if err := provider.AddNat(devIP, config.IP); err != nil {
		return fmt.Errorf("provider:add_nat: %s", err.Error())
	}

	return nil
}

// detachNetwork detaches the container from the host network
func detachNetwork(appModel *models.App, config docker.ContainerConfig) error {

	// fetch the devIP
	devIP := appModel.GlobalIPs["env"]

	//
	if err := provider.RemoveNat(devIP, config.IP); err != nil {
		return err
	}

	//
	if err := provider.RemoveIP(devIP); err != nil {
		return err
	}

	return nil
}

// devIsUnused returns true if the dev isn't being used by any other session
func devIsUnused() bool {
	count, err := counter.Get(config.EnvID() + "_dev")
	return count == 0 && err == nil
}

// isDevRunning returns true if a service is already running
func isDevRunning() bool {

	container, err := docker.GetContainer(container_generator.DevName())

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
