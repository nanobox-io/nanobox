// Package code ...
package code

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/validate"
)

// these constants represent different potential names a service can have
const (
	BUILD = "build"
)

// these constants represent different potential states an app can end up in
const (
	ACTIVE = "active"
)

func init() {
	validate.Register("built", validBuilt)
	validate.Register("dev_deployed", validDevDeployed)
	validate.Register("sim_deployed", validSimDeployed)
}

func validBuilt() error {
	env, err := models.FindEnvByID(config.EnvID())
	if err != nil || env.BuiltBoxfile == "" {
		return fmt.Errorf("No build has been completed for this application")
	}
	return nil
}

func validDevDeployed() error {
	app, err := models.FindAppBySlug(config.EnvID(), "dev")
	if err != nil || app.DeployedBoxfile == "" {
		return fmt.Errorf("Deploy has not been run for this application environment")
	}
	return nil
}

func validSimDeployed() error {
	app, err := models.FindAppBySlug(config.EnvID(), "sim")
	if err != nil || app.DeployedBoxfile == "" {
		return fmt.Errorf("Deploy has not been run for this application environment")
	}
	return nil
}
