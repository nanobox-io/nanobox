package dev

import (
	"github.com/spf13/cobra"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// StopCmd ...
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops your dev platform.",
		Long: `
Stops your dev platform. All data will be preserved in its current state.
		`,
		PreRun: steps.Run("start", "dev start"),
		Run:    stopFn,
	}
)

func init() {
	steps.Build("dev stop", true, stopCheck, stopFn)
}

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	// TODO: check the app and return some message
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")

	display.CommandErr(dev.Stop(app))
}

func stopCheck() bool {
	container, err := docker.GetContainer(container_generator.DevName())

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}