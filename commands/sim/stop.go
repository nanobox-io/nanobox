package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StopCmd ...
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops your sim platform.",
		Long: `
Stops your sim platform. All data and code
will be preserved in its current state.
		`,
		PreRun: validate.Requires("provider"),
		Run:    stopFn,
	}
)

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	simStop := sim.Stop{
		App: app,
	}
	display.CommandErr(simStop.Run())
}
