// Package code ...
package code

import (
	"errors"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
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

var errMissingImageOrName = errors.New("missing image or name")

func init() {
	validate.Register("built", validBuilt)
	validate.Register("dev_deployed", validDevDeployed)
	validate.Register("sim_deployed", validSimDeployed)
}

func validBuilt() error {
	boxfile := models.Boxfile{}
	if err := data.Get(config.AppID()+"_meta", "build_boxfile", &boxfile); err != nil {
		return errors.New("No build has been completed for this application")
	}
	return nil
}

func validDevDeployed() error {
	boxfile := models.Boxfile{}
	if err := data.Get(config.AppID()+"_meta", "dev_build_boxfile", &boxfile); err != nil {
		return errors.New("Deploy has not been run for this application environment")
	}
	return nil
}

func validSimDeployed() error {
	boxfile := models.Boxfile{}
	if err := data.Get(config.AppID()+"_meta", "sim_build_boxfile", &boxfile); err != nil {
		return errors.New("Deploy has not been run for this application environment")
	}
	return nil
}
