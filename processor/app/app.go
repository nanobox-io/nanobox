// Package app ...
package app

import (
	"fmt"

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
	return nil
}

func simIsUp() error {
	app := models.App{}
	data.Get("apps", config.AppName()+"_dev", &app)
	if !(app.Status == UP) {
		return fmt.Errorf("the environment has not been started. Please run the start command")
	}
	return nil
}
