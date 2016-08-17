// Package app ...
package app

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/validate"
)

// these constants represent different potential states an app can end up in
const (
	// states
	ACTIVE      = "active"
	INITIALIZED = "initialized"

	// statuses
	UP   = "up"
	DOWN = "down"
)

func init() {
	validate.Register("dev_isup", validateDevUp)
	validate.Register("sim_isup", validatSimUp)
}

func validateDevUp() error {
	// initialize the docker system so i can communicate
	dockerInit()

	// get the app and make sure its up
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	if app.Status != UP {
		return fmt.Errorf("the environment has not been started. Please run the start command")
	}

	// ensure all the services are running
	if !appServicesRunning(app) {
		return fmt.Errorf("The app is running but some services are not. Try running stop then start to clean this anomaly. if the problem persists please contact nanobox")
	}
	return nil
}

func validatSimUp() error {
	// initialize the docker system so i can communicate
	dockerInit()

	// get the app and make sure its up
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	if !(app.Status == UP) {
		return fmt.Errorf("the environment has not been started. Please run the start command")
	}

	// ensure all the services are running
	if !appServicesRunning(app) {
		return fmt.Errorf("The app is running but some services are not. Try running stop then start to clean this anomaly. if the problem persists please contact nanobox")
	}

	return nil
}

func appServicesRunning(app models.App) bool {

	components, _ := models.AllComponentsByApp(app.ID)
	for _, component := range components {
		// if any service are not running the app is not up
		if !isServiceRunning(component) {
			return false
		}
	}

	// all app services are running
	return true
}

// isServiceRunning returns true if a service is already running
func isServiceRunning(component models.Component) bool {

	// get the container
	container, err := docker.GetContainer(component.ID)

	if err != nil {
		lumber.Error("app:isServiceRunning(%#v): %s", component, err.Error())
		return false
	}

	return container.State.Status == "running"
}

func dockerInit() {
	provider.DockerEnv()
	docker.Initialize("env")
}
