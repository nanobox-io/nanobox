package processors

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	process_provider "github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/update"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/display"

)

func Update() error {

	// init docker client
	if err := process_provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	// check to see if nanobox needs to update
	update.Check()

	// pull the latest docker-machine image
	pullImages()

	// update all the nanobox images
	return provider.Install()
}

func pullImages() error{
	images, err := docker.ImageList()
	if err != nil {
		return err
	}

	for _, image := range images {
		display.StartTask("Pulling %s image", image.Slug)

		// generate a docker percent display
		dockerPercent := &display.DockerPercentDisplay{
			Output: display.NewStreamer("info"),
			// Prefix: buildImage,
		}

		// pull the build image
		if _, err := docker.ImagePull(image.Slug, dockerPercent); err != nil {
			lumber.Error("code:pullBuildImage:docker.ImagePull(%s, nil): %s", image.Slug, err.Error())
			display.ErrorTask()
			return fmt.Errorf("failed to pull docker image (%s): %s", image.Slug, err.Error())
		}
		display.StopTask()
	}
	return nil
}