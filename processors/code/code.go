// Package code ...
package code

import (
	"fmt"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/util"
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

func pullBuildImage() (string, error) {
	// extract the build image from the boxfile
	buildImage := buildImage()

	if docker.ImageExists(buildImage) {
		return buildImage, nil
	}

	display.StartTask("Pulling %s image", buildImage)
	defer display.StopTask()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: buildImage,
	}

	// pull the build image
	imagePull := func() error {
		_, err := docker.ImagePull(buildImage, dockerPercent)
		return err
	}
	if err := util.Retry(imagePull, 5, time.Second); err != nil {
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
		image = "nanobox/build"
	}

	return image
}
