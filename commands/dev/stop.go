package dev

import (
	"github.com/spf13/cobra"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

func init() {
	steps.Build("dev stop", true, stopCheck, stopFn)
}

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	// TODO: check the app and return some message
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	display.CommandErr(app.Stop(appModel))
}

func stopCheck() bool {
	container, err := docker.GetContainer(container_generator.DevName())

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}