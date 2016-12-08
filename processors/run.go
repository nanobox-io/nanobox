package processors

import (
	"fmt"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	build_generator "github.com/nanobox-io/nanobox/generators/hooks/build"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/watch"
)

// Run a code container with your runtime installed
func Run(envModel *models.Env, appModel *models.App, consoleConfig console.ConsoleConfig) error {

	// ensure the environment is setup
	if err := env.Setup(envModel); err != nil {
		return fmt.Errorf("failed to setup environment: %s", err.Error())
	}

	// setup the dev container
	if err := setup(appModel); err != nil {
		return fmt.Errorf("failed to setup dev container: %s", err.Error())
	}

	// start a watcher to watch for changes and inform the vm
	watchFiles(envModel, appModel)

	// create a dummy component using the appname
	component := &models.Component{
		ID: "nanobox_" + appModel.ID,
	}

	consoleConfig.DevIP = appModel.LocalIPs["env"]
	consoleConfig.Cwd = cwd(appModel)

	if err := env.Console(component, consoleConfig); err != nil {
		return fmt.Errorf("failed to console into dev container: %s", err)
	}

	if err := teardown(appModel); err != nil {
		return fmt.Errorf("unable to teardown dev: %s", err)
	}

	return nil
}

// sets up the dev container and network stack
func setup(appModel *models.App) error {

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// we don't need to setup if dev is already running
	if isDevExists() {
		if devInUse(container_generator.DevName()) {
			return nil
		} else {
			// if it isnt being used lets wipe it out and start again
			teardown(appModel)
		}
	}

	display.OpenContext("Building dev environment")
	defer display.CloseContext()

	// generate a container config
	config := container_generator.DevConfig(appModel)

	//
	if err := downloadImage(config.Image); err != nil {
		return err
	}

	display.StartTask("Starting docker container")
	container, err := docker.CreateContainer(config)
	if err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to create docker container: %s", err.Error())
	}
	display.StopTask()

	lumber.Prefix("dev:Console")
	defer lumber.Prefix("")

	display.StartTask("Configuring")
	userPayload := build_generator.UserPayload()
	if _, err := hookit.DebugExec(container.ID, "user", userPayload, "debug"); err != nil {
		return fmt.Errorf("failed to run the user hook: %s", err.Error())
	}

	if _, err := hookit.DebugExec(container.ID, "dev", build_generator.DevPayload(appModel), "info"); err != nil {
		return fmt.Errorf("failed to run the dev hook: %s", err.Error())
	}
	display.StopTask()

	return nil
}

func teardown(appModel *models.App) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	if devInUse(container_generator.DevName()) {
		return nil
	}

	// grab the container info
	container, err := docker.GetContainer(container_generator.DevName())
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}

	// remove the container
	if err := docker.ContainerRemove(container.ID); err != nil {
		lumber.Error("dev:console:teardown:docker.ContainerRemove(%s): %s", container.ID, err)
		return fmt.Errorf("failed to remove dev container: %s", err.Error())
	}

	return nil
}

// downloadImage downloads the dev docker image
func downloadImage(image string) error {

	if docker.ImageExists(image) {
		return nil
	}

	display.StartTask("Pulling %s image", image)
	defer display.StopTask()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: image,
	}

	imagePull := func() error {
		_, err := docker.ImagePull(image, dockerPercent)
		return err
	}
	if err := util.Retry(imagePull, 5, time.Second); err != nil {
		display.ErrorTask()
		lumber.Error("dev:Setup:downloadImage.ImagePull(%s, nil): %s", image, err.Error())
		return fmt.Errorf("failed to pull docker image (%s): %s", image, err.Error())
	}

	return nil
}

func watchFiles(envModel *models.Env, appModel *models.App) {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	if boxfile.Node("run.config").BoolValue("fs_watch") {
		lumber.Info("watcher starting")
		go watch.Watch(container_generator.DevName(), envModel.Directory)
	}
}

// cwd sets the cwd from the boxfile or provides a sensible default
func cwd(appModel *models.App) string {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	// parse the boxfile data

	if boxfile.Node("run.config").StringValue("cwd") != "" {
		return boxfile.Node("run.config").StringValue("cwd")
	}

	return "/app"
}

// devInUse returns true if the dev isn't being used by any other session
func devInUse(ID string) bool {
	consoles, _ := models.AllConsoles()
	for _, console := range consoles {
		// continue if the console isnt mine
		if console.ContainerID != ID {
			continue
		}
		if console.ID == "run" {
			return true
		}

		// check to see if this one is still running
		exec, err := docker.ExecInspect(console.ID)
		if err == nil && exec.Running {
			return true
		}

		// if we have a crufty exec delete it
		if err != nil || (err == nil && !exec.Running) {
			console.Delete()
		}
	}
	return false
}

// isDevExists returns true if a service is already running
func isDevExists() bool {

	_, err := docker.GetContainer(container_generator.DevName())

	// if the container doesn't exist then just return false
	return err == nil
}
