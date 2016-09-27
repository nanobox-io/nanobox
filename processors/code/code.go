// Package code ...
package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// these constants represent different potential names a service can have
const (
	BUILD = "build"
)

// these constants represent different potential states an app can end up in
const (
	ACTIVE = "active"
)

// runDebugSession drops the user in the build container to debug
func runDebugSession(container string, err error) error {
	if err != nil {
		display.ErrorTask()

		if registry.GetBool("debug") {
			component := &models.Component{
				ID: container,
			}
			err := env.Console(component, env.ConsoleConfig{})
			if err != nil {
				return fmt.Errorf("failed to establish a debug session: %s", err.Error())
			}
		}

	}

	return err
}

func pullBuildImage() (string, error) {
	// extract the build image from the boxfile
	buildImage := buildImage()

	display.StartTask("Pulling %s image", buildImage)
	defer display.StopTask()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: buildImage,
	}

	// pull the build image
	if _, err := docker.ImagePull(buildImage, dockerPercent); err != nil {
		lumber.Error("code:pullBuildImage:docker.ImagePull(%s, nil): %s", buildImage, err.Error())
		display.ErrorTask()
		return "", fmt.Errorf("failed to pull docker image (%s): %s", buildImage, err.Error())
	}

	return buildImage, nil
}

// BuildImage fetches the build image from the boxfile
func buildImage() string {
	// first let's see if the user has a custom build image they want to use
	box := boxfile.NewFromPath(config.Boxfile())
	image := box.Node("code.build").StringValue("image")

	// then let's set the default if the user hasn't specified
	if image == "" {
		image = "nanobox/build:proto"
	}

	return image
}
