package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
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
		PreRun: steps.Run("start", "sim start"),
		Run:    stopFn,
	}
)

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	display.CommandErr(sim.Stop(app))
}
