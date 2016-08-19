package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StopCmd ...
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops your dev platform.",
		Long: `
Stops your dev platform. All data will be preserved in its current state.
		`,
		PreRun: validate.Requires("provider"),
		Run:    stopFn,
	}
)

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	devStop := dev.Stop{
		App: app,
	}
	display.CommandErr(devStop.Run())
}
