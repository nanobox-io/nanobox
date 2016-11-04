package processors

import (
	"fmt"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	process_provider "github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
	//	"github.com/nanobox-io/nanobox/util/update"
)

func Update() error {

	// init docker client
	if err := process_provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	// // check to see if nanobox needs to update
	// update.Check()

	// update all the nanobox images
	pullImages()

	// pull the latest docker-machine image
	return provider.Install()
}

func pullImages() error {
	display.OpenContext("Updating Images")
	defer display.CloseContext()

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
		imagePullFunc := func() error {
			_, err := docker.ImagePull(image.Slug, dockerPercent)
			return err
		}

		if err := util.Retry(imagePullFunc, 5, time.Second); err != nil {
			lumber.Error("code:pullBuildImage:docker.ImagePull(%s, nil): %s", image.Slug, err.Error())
			display.ErrorTask()
			return fmt.Errorf("failed to pull docker image (%s): %s", image.Slug, err.Error())
		}

		display.StopTask()
	}

	return nil
}
