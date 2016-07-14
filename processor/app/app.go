// Package app ...
package app

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/validate"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
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
	validate.Register("dev_isup", devIsUp)
	validate.Register("sim_isup", simIsUp)
}

func devIsUp() error {
	app := models.App{}
	data.Get("apps", config.AppName()+"_dev", &app)
	if !(app.Status == UP) {
		return fmt.Errorf("the environment has not been started. Please run the start command")
	}
	if appServicesRunning(config.AppName()+"_dev") {
		return fmt.Errorf("the app is not up but some services are still running. Try running stop then start to clean this anomaly. if the problem persists please contact nanobox")
	}
	return nil
}

func simIsUp() error {
	app := models.App{}
	data.Get("apps", config.AppName()+"_dev", &app)
	if !(app.Status == UP) {
		return fmt.Errorf("the environment has not been started. Please run the start command")
	}

	if appServicesRunning(config.AppName()+"_dev") {
		return fmt.Errorf("the app is not up but some services are still running. Try running stop then start to clean this anomaly. if the problem persists please contact nanobox")
	}

	return nil
}


func appServicesRunning(app string) bool {

	// if the app cant be found in the database
	// its up and we will accept a failure later
	serviceNames, err := data.Keys(app)
	if err != nil {
		// if i cant get the keys from the app. its safer to assume the app is
		// down then to assume its up.
		return false
	}

	for _, serviceName := range serviceNames {
		// if any service is running the app is up
		if isServiceRunning(app, serviceName) {
			return true
		}
	}

	return false
}

// isServiceRunning returns true if a service is already running
func isServiceRunning(app, name string) bool {

	// get the container
	container, err := docker.GetContainer(fmt.Sprintf("nanobox_%s_%s", app, name))

	if err != nil {
		lumber.Error("Tried looking up nanobox_%s_%s Error: %s", app, name, err.Error())
		return false
	}

	return container.State.Status == "running"
}



