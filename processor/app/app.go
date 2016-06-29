// Package app ...
package app

import (
	"github.com/nanobox-io/nanobox/models"
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

func isUp() bool {
	app := models.App{}
	data.Get("apps", config.AppName()+"_dev", &app)
	return app.Status == UP
}
