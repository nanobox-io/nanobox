// Package code ...
package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/boxfile"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
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

// runDebugSession drops the user in the build container to debug
func runDebugSession(container string, err error) error {
	if registry.GetBool("debug") && err != nil {
		component := &models.Component{
			ID: container,
		}
		err := env.Console(component, env.ConsoleConfig{})
		if err != nil {
			return fmt.Errorf("failed to establish a debug session: %s", err.Error())
		}
	}

	return err
}

func pullBuildImage() (string, error) {
	// extract the build image from the boxfile
	buildImage := boxfile.BuildImage()

	display.StartTask("Pulling %s image", buildImage)
	defer display.StopTask()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		Prefix: buildImage,
	}

	// pull the build image
	if _, err := docker.ImagePull(buildImage, dockerPercent); err != nil {
		lumber.Error("code:pullBuildImage:docker.ImagePull(%s, nil): %s", buildImage, err.Error())
		display.ErrorTask()
		return "", fmt.Errorf("failed to pull docker image (%s): %s", buildImage, err.Error())
	}

	return buildImage, nil
}
